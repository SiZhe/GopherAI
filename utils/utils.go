package utils

import (
	"GopherAI/common/mysql/model"
	"regexp"

	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

// 获得随机数
func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < num; i++ {
		// 0~9随机数
		digit := r.Intn(10)
		code += strconv.Itoa(digit)
	}
	return code
}

// 生成唯一id
func GenerateUUID() string {
	return uuid.New().String()
}

// MD5 MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 将 schema 消息转换为数据库可存储的格式
func ConvertToModelMessage(sessionID string, userName string, msg *schema.Message) *model.Message {
	return &model.Message{
		SessionID: sessionID,
		UserName:  userName,
		Content:   msg.Content,
	}
}

// 将数据库消息转换为 schema 消息（供 AI 使用）
func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsgs := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		role := schema.Assistant
		if m.IsUser {
			role = schema.User
		}
		schemaMsgs = append(schemaMsgs, &schema.Message{
			Role:    role,
			Content: m.Content,
		})
	}
	return schemaMsgs
}

// ValidateFile 校验文件是否为允许的文本文件（.md）
func ValidateFile(file *multipart.FileHeader) error {
	// 校验文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".md" {
		return fmt.Errorf("文件类型不正确，只允许 .md，当前扩展名: %s", ext)
	}
	return nil
}

// 清除milvus表名称的违规符号
func CleanViolateSymbols(s string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(s, "")
}

// 返回：IP, 浏览器(含版本号)
func ParseIPAndBrowser(deviceInfo string) (string, string) {
	// 1. 提取IP
	ipRegex := regexp.MustCompile(`^(?:\d{1,3}\.){3}\d{1,3}`)
	ip := ipRegex.FindString(deviceInfo)
	if ip == "" {
		return "", "unknown"
	}

	// 去掉IP，剩下的是UA
	ua := strings.TrimPrefix(deviceInfo, ip)

	var browser, version string

	// 解析浏览器 + 版本
	switch {
	case regexp.MustCompile(`Version\/(\d+\.\d+)`).MatchString(ua):
		browser = "Safari"
		version = regexp.MustCompile(`Version\/(\d+\.\d+)`).FindStringSubmatch(ua)[1]

	case regexp.MustCompile(`Chrome\/(\d+\.\d+)`).MatchString(ua):
		browser = "Chrome"
		version = regexp.MustCompile(`Chrome\/(\d+\.\d+)`).FindStringSubmatch(ua)[1]

	case regexp.MustCompile(`Edg\/(\d+\.\d+)`).MatchString(ua):
		browser = "Edge"
		version = regexp.MustCompile(`Edg\/(\d+\.\d+)`).FindStringSubmatch(ua)[1]

	case regexp.MustCompile(`Firefox\/(\d+\.\d+)`).MatchString(ua):
		browser = "Firefox"
		version = regexp.MustCompile(`Firefox\/(\d+\.\d+)`).FindStringSubmatch(ua)[1]

	default:
		return ip, "unknown"
	}

	// 把浏览器 + 版本拼在一起返回
	return ip, browser + " " + version
}
