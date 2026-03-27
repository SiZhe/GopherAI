package rag

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

func EmbeddingModel(ctx context.Context) *openai.Embedder {
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		BaseURL: os.Getenv("SILICON_URL"),
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		Model:   os.Getenv("QWEN_EMBEDDING"),
	})
	if err != nil {
		panic(err)
	}
	return embedder
}
