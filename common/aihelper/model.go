package aihelper

import (
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
	DOUBAO_SEED_V16      string = "DOUBAO_SEED_V16"
	DOUBAO_SEED_V16_LITE string = "DOUBAO_SEED_V16_LITE"
	DOUBAO_SEED_CODE     string = "DOUBAO_SEED_CODE"
	DEEPSEEK_V32         string = "DEEPSEEK_V32"
)

type StreamCallback func(msg string)

// AIModel 定义AI模型接口
type AIModel interface {
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
	GetModelType() string
}

// =================== ARK实现 ===================
type ARKAIModel struct {
	llm       model.ToolCallingChatModel
	modelType string
}

// 根据modelType生成模型
func NewARKModel(ctx context.Context, modelType string) (*ARKAIModel, error) {
	key := os.Getenv("ARK_API_KEY")

	targetModelName := os.Getenv("ARK_MODEL_NAME_" + modelType)
	fmt.Println("modelName:", targetModelName)

	llm, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: key,
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
