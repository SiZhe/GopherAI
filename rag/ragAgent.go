package rag

/*
import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"GopherAI/common/aihelper"
)

type RagModel struct {
	model    aihelper.RagAIModel
	messages []*schema.Message
}

var RagAgent *compose.Graph[RagModel, []*schema.Message]

type Judge struct {
	model    aihelper.RagAIModel
	messages []*schema.Message
	flag     bool
}

func aaa() {
	//ctx := context.Background()
	graph := compose.NewGraph[[]*schema.Message, []*schema.Message]()

	// 节点1
	lambdaJudge := compose.InvokableLambda(func(ctx context.Context, input RagModel) (output Judge, err error) {
		isExistFiles, err := IsExistUploadsFiles(input.model.GetUserName(), input.model.GetSessionId())
		if err != nil {
			return Judge{}, err
		}
		fmt.Printf("isExistFiles:%v\n\n", isExistFiles)
		// 不存在直接走普通model
		if !isExistFiles {
			return Judge{
				model:    input.model,
				messages: input.messages,
				flag:     false,
			}, nil
		}

		// 有文件->判断要不要走rag
		// 1.取最后一条消息
		query := input.messages[len(input.messages)-1]
		judgeQuery := BuildJudgePrompt(query)
		fmt.Printf("query:%v\n\n", query.Content)

		// 2.输入给小型llm 判断要不要走rag
		llm, err := aihelper.NewRagJudgeArkModel("DOUBAO_15_LITE")
		if err != nil {
			return Judge{}, err
		}

		judge, err := llm.Generate(ctx, judgeQuery)
		if err != nil {
			return Judge{}, err
		}
		if judge.Content == "true" {
			return Judge{
				model:    input.model,
				messages: input.messages,
				flag:     true,
			}, nil
		} else if judge.Content == "false" {
			return Judge{
				model:    input.model,
				messages: input.messages,
				flag:     false,
			}, nil
		} else {
			return Judge{}, errors.New("模型输出错误，非预期的 true/false")
		}
	})

}

func BuildJudgePrompt(query *schema.Message) []*schema.Message {
	// 构造系统提示词（你的分类规则）
	systemContent := `# 任务：RAG路由二分类器
你是一个严格的二分类器，只负责判断用户问题是否需要查询知识库。

## 绝对规则
1. 需要走 RAG（返回 true）：
- 任何涉及公司内部制度、流程、业务数据、产品文档、项目资料的问题
- 任何涉及遥感、GIS、地理信息、NDVI、NPP、生态遥感、卫星数据的专业问题
- 任何需要精确数据、最新信息、特定文档内容的问题
- 任何你不确定、不知道、需要查证的问题

2. 不需要走 RAG（返回 false）：
- 所有闲聊、问候、打招呼、感谢、告别
- 所有通用常识、简单逻辑、数学计算
- 所有大模型本身就能准确回答的通用知识问题
- 所有关于你自己的问题

## 输出要求（违反即错误）
1. 绝对只输出一个单词：true 或者 false
2. 不要任何解释、不要任何理由、不要任何标点
3. 不要大写、不要引号、不要JSON、不要任何其他字符
4. 不要说"我认为"、"应该"、"可能"这类不确定的词

## 示例
用户：你好 → 输出：false
用户：今天天气怎么样 → 输出：true
用户：1+1等于几 → 输出：false
用户：什么是NPP → 输出：true
用户：公司请假流程是什么 → 输出：true
用户：帮我写一段代码 → 输出：false
用户：地球是圆的吗 → 输出：false
用户：吉林一号卫星分辨率是多少 → 输出：true`

	// 返回标准 messages 结构：系统消息 + 用户消息
	return []*schema.Message{
		schema.SystemMessage(systemContent),
		schema.UserMessage(query.Content),
	}
}
*/
