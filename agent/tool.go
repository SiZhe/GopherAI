package agent

import (
	"GopherAI/config"
	"GopherAI/rag"
	"os"

	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// ==================== CheckDocTool 定义 ====================

// CheckDocInputParams 检查文档的输入参数
type CheckDocInputParams struct {
	UserName  string `json:"user_name" jsonschema:"description=用户名，用于查询该用户的文档上传情况"`
	SessionID string `json:"session_id" jsonschema:"description=会话ID，用于查询该会话的文档上传情况"`
}

// CheckDocOutputParams 检查文档的输出结果
type CheckDocOutputParams struct {
	HasDoc  bool   `json:"has_doc" jsonschema:"description=是否有上传的文档，true表示有，false表示没有"`
	Message string `json:"message,omitempty" jsonschema:"description=结果描述信息"`
}

// CheckDocTool 检查当前会话是否有上传文档的工具
func CheckDocTool(_ context.Context, input *CheckDocInputParams) (*CheckDocOutputParams, error) {
	fmt.Printf("[CheckDocTool] 开始检查文档, username: %s, sessionID: %s\n", input.UserName, input.SessionID)

	// 参数验证
	if input.UserName == "" {
		return &CheckDocOutputParams{
			HasDoc:  false,
			Message: "用户名不能为空",
		}, fmt.Errorf("username is empty")
	}
	if input.SessionID == "" {
		return &CheckDocOutputParams{
			HasDoc:  false,
			Message: "会话ID不能为空",
		}, fmt.Errorf("session_id is empty")
	}

	// 调用 rag 包检查文档
	hasDoc, err := rag.IsExistUploadsFiles(input.UserName, input.SessionID)
	if err != nil {
		fmt.Printf("[CheckDocTool] 检查文档失败: %v\n", err)
		return &CheckDocOutputParams{
			HasDoc:  false,
			Message: fmt.Sprintf("检查文档失败: %v", err),
		}, err
	}

	// 构建结果
	message := ""
	if hasDoc {
		message = fmt.Sprintf("会话 %s 有上传的文档", input.SessionID)
	} else {
		message = fmt.Sprintf("会话 %s 没有上传的文档", input.SessionID)
	}

	fmt.Printf("[CheckDocTool] 检查完成: hasDoc=%v, message=%s\n", hasDoc, message)

	return &CheckDocOutputParams{
		HasDoc:  hasDoc,
		Message: message,
	}, nil
}

// CreateCheckDocTool 创建检查文档的工具（修复拼写错误）
func CreateCheckDocTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "check_documents",
			Desc: "检查用户在当前会话中是否有上传的文档，用于决定是否需要使用 RAG 检索",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"user_name": {
					Type:     schema.String,
					Desc:     "用户名，用于查询该用户的文档上传情况，可从system messages中获得",
					Required: true,
				},
				"session_id": {
					Type:     schema.String,
					Desc:     "会话ID，用于查询该会话的文档上传情况，可从system messages中获得",
					Required: true,
				},
			}),
		},
		CheckDocTool,
	)
}

// ==================== RAGDocTool 定义 ====================

// RAGInputParams RAG 检索的输入参数
type RAGInputParams struct {
	UserName  string `json:"user_name" jsonschema:"description=用户名，用于查询该用户的文档"`
	SessionID string `json:"session_id" jsonschema:"description=会话ID，用于查询该会话的文档"`
	Query     string `json:"query" jsonschema:"description=用户的最后一个问题！必须传入用户最近一次提问的内容，不能修改或简化，用于文档检索"`
}

// RAGOutputParams RAG 检索的输出结果
type RAGOutputParams struct {
	Success   bool               `json:"success" jsonschema:"description=检索是否成功"`
	Message   string             `json:"message,omitempty" jsonschema:"description=结果描述信息"`
	DocsCount int                `json:"docs_count" jsonschema:"description=检索到的文档数量"`
	Documents []*schema.Document `json:"documents,omitempty" jsonschema:"description=检索到的文档列表"`
	RAGResult string             `json:"rag_result,omitempty" jsonschema:"description=检索到的 RAG 文档"`
}

// RAGDocTool 执行 RAG 文档检索的工具
func RAGTool(_ context.Context, input *RAGInputParams) (*RAGOutputParams, error) {
	// 参数验证
	if input.UserName == "" {
		return &RAGOutputParams{
			Success: false,
			Message: "用户名不能为空",
		}, fmt.Errorf("username is empty")
	}
	if input.SessionID == "" {
		return &RAGOutputParams{
			Success: false,
			Message: "会话ID不能为空",
		}, fmt.Errorf("session_id is empty")
	}
	if input.Query == "" {
		return &RAGOutputParams{
			Success: false,
			Message: "用户问题不能为空，必须传入用户的最后一个问题",
		}, fmt.Errorf("query is empty")
	}

	fmt.Printf("[RAGDocTool] 开始 RAG 检索, username: %s, sessionID: %s, query: %s\n",
		input.UserName, input.SessionID, input.Query)

	// 创建检索器
	ctx := context.Background()
	retriever := rag.NewRAGRetriever(ctx, input.UserName, input.SessionID, config.GetConfig().RagConfig.TopK)

	// 执行检索
	fmt.Printf("[RAGDocTool] 执行检索, query: %s\n", input.Query)
	retrieveDocs, err := retriever.RetrieverUploadsFiles(ctx, input.Query)
	if err != nil {
		fmt.Printf("[RAGDocTool] 检索失败: %v\n", err)
		return &RAGOutputParams{
			Success: false,
			Message: fmt.Sprintf("检索失败: %v", err),
		}, err
	}

	for i, doc := range retrieveDocs {
		fmt.Printf("[RAGDocTool] 检索到文档%d: %s\n", i+1, doc.Content)
	}

	// 构建 RAG 提示词
	ragResult := rag.BuildRAGDocuments(retrieveDocs)
	fmt.Printf("[RAGDocTool] 构建 RAG 提示词完成, 长度: %d\n", len(ragResult))

	// 构建结果
	result := &RAGOutputParams{
		Success:   true,
		DocsCount: len(retrieveDocs),
		Documents: retrieveDocs,
		RAGResult: ragResult,
	}

	if len(retrieveDocs) > 0 {
		result.Message = fmt.Sprintf("成功检索到 %d 个相关文档", len(retrieveDocs))
	} else {
		result.Message = "未检索到相关文档"
	}

	fmt.Printf("[RAGDocTool] 检索完成: success=%v, docsCount=%d, message=%s\n",
		result.Success, result.DocsCount, result.Message)

	return result, nil
}

// CreateRAGTool 创建 RAG 文档检索工具
func CreateRAGTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "retrieve_documents",
			Desc: "根据用户的最后一个问题，从当前会话的上传文档中检索相关内容，用于 RAG 增强生成。重要：query 参数必须传入用户的最后一个问题，不能修改或简化！",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"user_name": {
					Type:     schema.String,
					Desc:     "用户名，用于查询该用户的文档",
					Required: true,
				},
				"session_id": {
					Type:     schema.String,
					Desc:     "会话ID，用于查询该会话的文档",
					Required: true,
				},
				"query": {
					Type:     schema.String,
					Desc:     "用户的最后一个问题！必须传入用户最近一次提问的内容，不能修改或简化，不能省略任何关键信息，用于文档检索",
					Required: true,
				},
			}),
		},
		RAGTool,
	)
}

// ==================== JudgeTool 定义 ====================

// JudgeInputParams RAG 决策的输入参数
type JudgeInputParams struct {
	Query     string `json:"query" jsonschema:"description=用户的问题，用于判断是否需要 RAG 检索"`
	RAGResult string `json:"rag_result,omitempty" jsonschema:"description=根据用户的问题所构建的 RAG 提示词"`
}

// JudgeOutputParams RAG 决策的输出结果
type JudgeOutputParams struct {
	NeedRAG bool   `json:"need_rag" jsonschema:"description=是否需要使用 RAG，true表示需要，false表示不需要"`
	Message string `json:"message,omitempty" jsonschema:"description=判断理由"`
	Prompt  string `json:"prompt,omitempty" jsonschema:"description=结果描述信息"`
}

// JudgeTool 使用小型模型判断是否需要 RAG 的工具
func JudgeTool(ctx context.Context, input *JudgeInputParams) (*JudgeOutputParams, error) {
	fmt.Printf("[JudgeTool] 开始 RAG 决策, query: %s\n", input.Query)

	// 参数验证
	if input.Query == "" {
		return &JudgeOutputParams{
			NeedRAG: false,
			Message: "用户问题不能为空",
			Prompt:  "",
		}, fmt.Errorf("query is empty")
	}

	llm, err := newJudgeArkModel("DOUBAO_15_LITE")
	if err != nil {
		fmt.Printf("[JudgeTool] 初始化模型失败: %v", err)
		// 降级：使用关键词判断
		return keywordRAGJudge(input.Query), nil
	}

	// 构建决策提示词
	judgeMessages := buildJudgePrompt(input.Query, input.RAGResult)

	// 调用模型进行决策
	fmt.Printf("[JudgeTool] 调用模型进行决策\n")
	judge, err := llm.Generate(ctx, judgeMessages)
	if err != nil {
		fmt.Printf("[RAGJudgeTool] 模型调用失败: %v, 降级为关键词判断\n", err)
		// 降级：使用关键词判断
		return keywordRAGJudge(input.Query), nil
	}

	// 模型只会输出true或者false
	if judge.Content == "true" {
		return &JudgeOutputParams{
			NeedRAG: true,
			Message: "模型输出模糊匹配为需要 RAG",
			Prompt:  rag.BuildRAGPromptString(input.Query, input.RAGResult),
		}, nil
	} else if judge.Content == "false" {
		return &JudgeOutputParams{
			NeedRAG: false,
			Message: "模型判断不需要使用 RAG",
			Prompt:  input.Query,
		}, nil
	} else {
		return &JudgeOutputParams{}, fmt.Errorf("模型输出错误:%s，非预期 true/false", judge.Content)
	}
}

// buildJudgePrompt 构建 RAG 决策提示词
func buildJudgePrompt(query string, ragResult string) []*schema.Message {
	systemMessage := `# 任务：RAG路由二分类器
你是一个严格的二分类器，根据【用户问题】和【RAG检索结果】判断是否需要使用 RAG。

## 判断规则
1. 返回 true（使用 RAG）：
- 检索到的文档中包含回答用户问题的相关内容
- 用户问题涉及文档中的特定信息、数据、概念
- 用户问题明确要求查询文档内容
- 你认为文档内容对回答问题有帮助

2. 返回 false（不使用 RAG）：
- 没有检索到任何相关文档
- 检索到的文档与用户问题完全无关
- 用户问题是简单闲聊、问候、通用常识
- 你认为文档内容对回答问题没有帮助

## 输出要求（严格遵守）
1. 绝对只输出一个单词：true 或者 false
2. 不要任何解释、不要任何理由、不要任何标点
3. 不要大写、不要引号、不要JSON、不要任何其他字符`

	// 构建用户消息（包含问题和检索结果）
	userMessage := "【用户问题】:" + query + "。" + "\n\n" + "【RAG检索结果】:" + ragResult + "\n\n"

	return []*schema.Message{
		schema.SystemMessage(systemMessage),
		schema.UserMessage(userMessage),
	}
}

// keywordRAGJudge 关键词判断（降级方案）
func keywordRAGJudge(query string) *JudgeOutputParams {
	queryLower := strings.ToLower(query)

	// 不需要 RAG 的关键词（简单对话）
	noRAGKeywords := []string{
		"你好", "哈喽", "hi", "hello", "嗨",
		"再见", "拜拜", "bye",
		"谢谢", "感谢",
		"好的", "嗯", "哦", "知道了",
		"1+", "等于", "多少",
		"地球", "太阳", "月亮",
	}

	for _, kw := range noRAGKeywords {
		if strings.Contains(queryLower, kw) && len(query) < 15 {
			return &JudgeOutputParams{
				NeedRAG: false,
				Message: "关键词判断：简单对话，不需要 RAG",
				Prompt:  query,
			}
		}
	}

	// 需要 RAG 的关键词
	ragKeywords := []string{
		"文档", "文件", "资料", "参考", "根据", "上传",
		"内容", "里面", "中", "里",
		"什么", "多少", "几", "谁", "哪", "哪里", "哪个",
		"解释", "说明", "介绍", "描述",
		"数据", "信息", "知识", "概念",
		"流程", "制度", "规定",
		"npp", "ndvi", "遥感", "gis", "地理",
	}

	for _, kw := range ragKeywords {
		if strings.Contains(queryLower, kw) {
			return &JudgeOutputParams{
				NeedRAG: true,
				Message: "关键词判断：检测到需要文档信息的关键词",
				Prompt:  query,
			}
		}
	}

	// 默认需要 RAG（保守策略）
	return &JudgeOutputParams{
		NeedRAG: true,
		Message: "关键词判断：默认需要 RAG",
		Prompt:  query,
	}
}

func newJudgeArkModel(modelType string) (model.ChatModel, error) {
	api_key := os.Getenv("ARK_API_KEY")

	targetModelName := ""

	if modelType == "DOUBAO_15_LITE" {
		targetModelName = os.Getenv("ARK_DOUBAO_15_LITE")
		fmt.Println("modelName:", targetModelName)
	} else {
		return nil, fmt.Errorf("invalid model name")
	}

	ctx := context.Background()
	llm, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: api_key,
		Model:  targetModelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create ark model failed: %v", err)
	}
	return llm, nil
}

// CreateJudgeTool 创建 RAG 决策工具
func CreateJudgeTool() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "judge",
			Desc: "根据用户问题和 RAG 检索到的文档内容，判断是否需要使用 RAG。",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"query": {
					Type:     schema.String,
					Desc:     "用户的最后一个问题！必须传入用户最近一次提问的内容，不能修改或简化，不能省略任何关键信息",
					Required: true,
				},
				"rag_result": {
					Type:     schema.String,
					Desc:     "根据用户的问题得到的 RAG 结果，可以调用 RAGTool 获取检索结果",
					Required: true,
				},
			}),
		},
		JudgeTool,
	)
}
