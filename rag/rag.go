package rag

import (
	"GopherAI/common/milvus"
	"GopherAI/utils"
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	milvusIndexer "github.com/cloudwego/eino-ext/components/indexer/milvus"
	milvusRetriever "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type RAGIndexer struct {
	embedding *openai.Embedder
	indexer   *milvusIndexer.Indexer
}

type RAGRetriever struct {
	embedding *openai.Embedder
	retriever *milvusRetriever.Retriever
}

func NewRAGIndexer(ctx context.Context, username string, sessionId string) *RAGIndexer {
	var fields = []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "256",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector",
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "81920",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}

	embeddingModel := EmbeddingModel(ctx)

	collection := utils.CleanViolateSymbols(fmt.Sprintf("user_%s_session_%s", username, sessionId))
	fmt.Println(collection)

	indexer, err := milvusIndexer.NewIndexer(ctx, &milvusIndexer.IndexerConfig{
		Client: milvus.GetMilvusClient(),
		// 为每个对话建表
		Collection: collection,
		Fields:     fields,
		Embedding:  embeddingModel,
	})
	if err != nil {
		panic(err)
	}

	return &RAGIndexer{
		embedding: embeddingModel,
		indexer:   indexer,
	}
}

func NewRAGRetriever(ctx context.Context, username string, sessionId string, topK int) *RAGRetriever {
	embeddingModel := EmbeddingModel(ctx)

	collection := utils.CleanViolateSymbols(fmt.Sprintf("user_%s_session_%s", username, sessionId))
	fmt.Println(collection)

	retrieve, err := milvusRetriever.NewRetriever(ctx, &milvusRetriever.RetrieverConfig{
		Client: milvus.GetMilvusClient(),
		// 每个对话一个表
		Collection:  collection,
		VectorField: "vector",
		OutputFields: []string{
			"id", "content", "metadata",
		},
		TopK:      topK,
		Embedding: embeddingModel,
	})
	if err != nil {
		panic(err)
	}

	return &RAGRetriever{
		embedding: embeddingModel,
		retriever: retrieve,
	}
}

func (ragIndexer *RAGIndexer) IndexerUploadsFiles(ctx context.Context, docs []*schema.Document) (ids []string, err error) {
	idx, err := ragIndexer.indexer.Store(ctx, docs)
	if err != nil {
		return []string{}, err
	} else {
		return idx, nil
	}
}

func (ragRetriever *RAGRetriever) RetrieverUploadsFiles(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	retrieveDocs, err := ragRetriever.retriever.Retrieve(ctx, query, opts...)
	if err != nil {
		return []*schema.Document{}, err
	} else {
		return retrieveDocs, nil
	}
}

// BuildRAGPrompt 构建包含检索文档的提示词
func BuildRAGPrompt(query string, docs []*schema.Document) string {
	if len(docs) == 0 {
		return query
	}

	contextText := ""
	for i, doc := range docs {
		contextText += fmt.Sprintf("[文档 %d]: %s\n\n", i+1, doc.Content)
	}

	prompt := fmt.Sprintf(`基于以下参考文档回答用户的问题。如果文档中没有相关信息，请自由回答。

参考文档：%s 

用户问题：%s

请提供准确、完整的回答：`, contextText, query)

	return prompt
}
