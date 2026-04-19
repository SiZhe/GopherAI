package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

type RagReActAgent struct {
	ReActAgent *react.Agent
}

func NewRagReActAgent(llm model.ToolCallingChatModel) (RagReActAgent, error) {
	ctx := context.Background()
	checkDocTool := CreateCheckDocTool()
	ragTool := CreateRAGTool()
	judgeTool := CreateJudgeTool()

	reactAgent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: llm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools:               []tool.BaseTool{checkDocTool, ragTool, judgeTool},
			ExecuteSequentially: false,
		},
		MaxStep:       12,
		GraphName:     "",
		ModelNodeName: "",
		ToolsNodeName: "",
	})

	if err != nil {
		return RagReActAgent{}, err
	}
	return RagReActAgent{
		ReActAgent: reactAgent,
	}, nil
}

func (ragAgent *RagReActAgent) Get(messages []*schema.Message) (*schema.Message, error) {
	inputMessages := make([]*schema.Message, 0)
	for _, m := range messages {
		if m.Role == schema.System {
			inputMessages = append(inputMessages, m)
		}
	}

	inputMessages = append(inputMessages, messages[len(messages)-1])

	inputMessages = append(inputMessages, schema.SystemMessage(
		"你是一个用户问题智能助手，根据工具，是否需要使用 RAG技术，如果需要就将用户的问题改成 RAG 修饰的提示词,如果不需要则直接返回原问题！"))
	ctx := context.Background()
	result, err := ragAgent.ReActAgent.Generate(ctx, inputMessages, agent.WithComposeOptions(compose.WithCallbacks(&loggerCallback{})))
	if err != nil {
		return nil, err
	}
	return result, nil
}

type loggerCallback struct {
	callbacks.HandlerBuilder
}

func (cb *loggerCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	fmt.Println("==================")
	inputStr, _ := json.MarshalIndent(input, "", "  ") // nolint: byted_s_returned_err_check
	fmt.Printf("[OnStart] %s\n", string(inputStr))
	return ctx
}

func (cb *loggerCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	fmt.Println("=========[OnEnd]=========")
	outputStr, _ := json.MarshalIndent(output, "", "  ") // nolint: byted_s_returned_err_check
	fmt.Println(string(outputStr))
	return ctx
}

func (cb *loggerCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	fmt.Println("=========[OnError]=========")
	fmt.Println(err)
	return ctx
}

func (cb *loggerCallback) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo,
	output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {

	var graphInfoName = react.GraphName

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[OnEndStream] panic err:", err)
			}
		}()

		defer output.Close() // remember to close the stream in defer

		fmt.Println("=========[OnEndStream]=========")
		for {
			frame, err := output.Recv()
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			s, err := json.Marshal(frame)
			if err != nil {
				fmt.Printf("internal error: %s\n", err)
				return
			}

			if info.Name == graphInfoName { // 仅打印 graph 的输出, 否则每个 stream 节点的输出都会打印一遍
				fmt.Printf("%s: %s\n", info.Name, string(s))
			}
		}

	}()
	return ctx
}

func (cb *loggerCallback) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo,
	input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	defer input.Close()
	return ctx
}
