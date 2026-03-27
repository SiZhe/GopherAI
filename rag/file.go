package rag

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"GopherAI/utils"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
)

func CreateUploadsDir(username string, sessionId string) error {
	// 用户文件目录
	userDir := filepath.Join("uploads_files", username)
	// 用户对话目录
	sessionDir := filepath.Join(userDir, sessionId)

	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		log.Printf("Failed to create user directory %s: %v", userDir, err)
		return err
	}
	return nil
}

func IsExistUploadsFiles(username string, sessionId string) (bool, error) {
	// 用户文件目录
	userDir := filepath.Join("uploads_files", username)
	// 用户对话目录
	sessionDir := filepath.Join(userDir, sessionId)

	//读取目录
	files, err := os.ReadDir(sessionDir)
	if err != nil {
		return false, err
	}

	return len(files) > 0, nil
}

// 返回文件路径
func StoreUploadsFiles(username string, sessionId string, file *multipart.FileHeader) (string, error) {
	// 校验文件类型和文件名
	if err := utils.ValidateFile(file); err != nil {
		log.Printf("File validation failed: %v", err)
		return "", err
	}

	// 文件名 ： 唯一id + 文件类型
	filename := utils.GenerateUUID() + filepath.Ext(file.Filename)

	filePath := filepath.Join("uploads_files", username, sessionId, filename)

	// 打开上传的文件
	srcFile, err := file.Open()
	if err != nil {
		log.Printf("Failed to open uploaded file: %v", err)
		return "", err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create destination file %s: %v", filePath, err)
		return "", err
	}
	defer dstFile.Close()

	// 写文件
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		log.Printf("Failed to copy file content: %v", err)
		return "", err
	} else {
		log.Printf("File uploaded successfully: %s", filePath)
	}

	return filePath, nil
}

func TransformerUploadsFiles(ctx context.Context, filePath string) ([]*schema.Document, error) {
	// 准备分割器
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		return nil, err
	}

	// 打开分割的文档
	content, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	defer content.Close()

	bs, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	//从路径得到用户名，对话id，文件名
	parts := filepath.SplitList(filePath)

	// 兼容 Windows 的 \ 分隔符，统一转成 /
	parts = strings.Split(filepath.ToSlash(filePath), "/")

	// 满足：uploads_files / 用户 / session / 文件
	if len(parts) < 4 || parts[0] != "uploads_files" {
		log.Println("invalid file path format")
	}

	username := parts[1]
	sessionId := parts[2]
	filename := parts[3]

	// fmt.Println(filename)

	doc := []*schema.Document{
		{
			// 这里的id没什么用
			ID:      utils.GenerateUUID(),
			Content: string(bs),
			MetaData: map[string]any{
				"username":  username,
				"sessionId": sessionId,
				"filename":  filename,
			},
		},
	}

	// 分割后的文档
	splitDocs, err := splitter.Transform(ctx, doc)

	for i, _ := range splitDocs {
		splitDocs[i].ID = utils.GenerateUUID()
	}

	if err != nil {
		return nil, err
	}

	return splitDocs, nil
}
