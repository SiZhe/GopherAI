package aihelper

import (
	"GopherAI/config"
	"GopherAI/rag"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	DOUBAO_SEED_20 string = "DOUBAO_SEED_20"
	DEEPSEEK_V32   string = "DEEPSEEK_V32"
)

type StreamCallback func(msg string)

// AIModel 定义AI模型接口
type RagAIModel interface {
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
	GetModelType() string
	GetUserName() string
	GetSessionId() string
}

// =================== ARK实现 ===================
/*
type ARKAIModel struct {
	llm       model.ToolCallingChatModel
	modelType string
}

// 根据modelType生成模型
func NewARKModel(ctx context.Context, modelType string) (*ARKAIModel, error) {
	api_key := os.Getenv("ARK_API_KEY")

	targetModelName := ""

	if modelType == DOUBAO_SEED_20 {
		targetModelName = os.Getenv("ARK_DOUBAO_SEED_20")
		fmt.Println("modelName:", targetModelName)
	} else if modelType == DEEPSEEK_V32 {
		targetModelName = os.Getenv("ARK_DEEPSEEK_V32")
		fmt.Println("modelName:", targetModelName)
	} else {
		return nil, fmt.Errorf("invalid model name")
	}

	llm, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: api_key,
		Model:  targetModelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create ark model failed: %v", err)
	}
	return &ARKAIModel{llm: llm, modelType: modelType}, nil
}

func (a *ARKAIModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := a.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("ark generate failed: %v", err)
	}
	return resp, nil
}

func (a *ARKAIModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := a.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("ark stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("ark stream recv failed: %v", err)
		}
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content) // 聚合
			cb(msg.Content)                   // 实时调用cb函数，方便主动发送给前端
		}
	}

	return fullResp.String(), nil //返回完整内容，方便后续存储
}

func (a *ARKAIModel) GetModelType() string {
	return a.modelType
}
*/

// =================== RAG-ARK实现 ===================
type RagArkAIModel struct {
	llm       model.ToolCallingChatModel
	modelType string
	username  string
	sessionId string
}

func NewRagArkModel(ctx context.Context, modelType string, username string, sessionId string) (*RagArkAIModel, error) {
	api_key := os.Getenv("ARK_API_KEY")

	targetModelName := ""

	if modelType == DOUBAO_SEED_20 {
		targetModelName = os.Getenv("ARK_DOUBAO_SEED_20")
		fmt.Println("modelName:", targetModelName)
	} else if modelType == DEEPSEEK_V32 {
		targetModelName = os.Getenv("ARK_DEEPSEEK_V32")
		fmt.Println("modelName:", targetModelName)
	} else {
		return nil, fmt.Errorf("invalid model name")
	}

	llm, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: api_key,
		Model:  targetModelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create ark model failed: %v", err)
	}
	return &RagArkAIModel{llm: llm, modelType: modelType, username: username, sessionId: sessionId}, nil
}

func (ragAIModel *RagArkAIModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	isExistFiles, err := rag.IsExistUploadsFiles(ragAIModel.username, ragAIModel.sessionId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("isExistFiles:%v\n\n", isExistFiles)

	// 有文件就rag，没有文件就不用
	if isExistFiles || len(messages) == 0 {
		fmt.Println("执行rag检索....\n")
		retriever := rag.NewRAGRetriever(ctx, ragAIModel.username, ragAIModel.sessionId, config.GetConfig().RagConfig.TopK)

		//取最后一条消息
		query := messages[len(messages)-1].Content
		fmt.Printf("query:%v\n\n", query)

		retrieveDocs, err := retriever.RetrieverUploadsFiles(ctx, query)
		for i, doc := range retrieveDocs {
			fmt.Printf("retrieveDocs[%v]:%v\n", i+1, doc.Content)
		}
		if err != nil {
			return nil, err
		}

		// 构建提示词
		ragPrompt := rag.BuildRAGPrompt(query, retrieveDocs)
		fmt.Printf("ragPrompt:%v\n\n", ragPrompt)

		// 替换最后一条消息为 RAG 提示词
		messages[len(messages)-1] = schema.UserMessage(ragPrompt)
	}

	resp, err := ragAIModel.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("ark generate failed: %v", err)
	}
	return resp, nil
}

func (ragAIModel *RagArkAIModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	isExistFiles, err := rag.IsExistUploadsFiles(ragAIModel.username, ragAIModel.sessionId)
	if err != nil {
		return "", err
	}
	fmt.Printf("isExistFiles:%v\n\n", isExistFiles)

	// 有文件就rag，没有文件就不用
	if isExistFiles || len(messages) == 0 {
		fmt.Println("执行rag检索....\n")
		retriever := rag.NewRAGRetriever(ctx, ragAIModel.username, ragAIModel.sessionId, config.GetConfig().RagConfig.TopK)

		//取最后一条消息
		query := messages[len(messages)-1].Content
		fmt.Printf("query:%v\n\n", query)

		retrieveDocs, err := retriever.RetrieverUploadsFiles(ctx, query)
		if err != nil {
			return "", err
		}

		// 构建提示词
		ragPrompt := rag.BuildRAGPrompt(query, retrieveDocs)
		fmt.Printf("ragPrompt:%v\n\n", ragPrompt)

		// 替换最后一条消息为 RAG 提示词
		messages[len(messages)-1] = schema.UserMessage(ragPrompt)
	}

	stream, err := ragAIModel.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("ark stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("ark stream recv failed: %v", err)
		}
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content) // 聚合
			cb(msg.Content)                   // 实时调用cb函数，方便主动发送给前端
		}
	}

	return fullResp.String(), nil //返回完整内容，方便后续存储
}

func (a *RagArkAIModel) GetModelType() string {
	return a.modelType
}

func (ragAIModel *RagArkAIModel) GetUserName() string {
	return ragAIModel.username
}

func (ragAIModel *RagArkAIModel) GetSessionId() string {
	return ragAIModel.sessionId
}
