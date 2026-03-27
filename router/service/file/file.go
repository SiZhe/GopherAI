package file

import (
	"GopherAI/rag"
	"context"
	"mime/multipart"
)

func UploadsFiles(ctx context.Context, username string, sessionId string, uploadedFile *multipart.FileHeader) (filePath string, err error) {
	// 创建目录(确保)
	err = rag.CreateUploadsDir(username, sessionId)
	if err != nil {
		return "", err
	}

	// 将文件保存到服务器
	filePath, err = rag.StoreUploadsFiles(username, sessionId, uploadedFile)
	if err != nil {
		return "", err
	}

	// 分割文件
	docs, err := rag.TransformerUploadsFiles(ctx, filePath)
	if err != nil {
		return "", err
	}

	// 得到indexer
	ragIndexer := rag.NewRAGIndexer(ctx, username, sessionId)
	// 存储到数据库中
	_, err = ragIndexer.IndexerUploadsFiles(ctx, docs)
	if err != nil {
		return "", err
	}
	return filePath, err
}
