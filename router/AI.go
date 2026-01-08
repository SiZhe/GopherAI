package router

import (
	"GopherAI/router/controller/session"

	"github.com/gin-gonic/gin"
)

func AIRouter(r *gin.RouterGroup) {
	// 聊天相关接口
	r.GET("/chat/sessions", session.GetUserSessionsByUserName)
	r.POST("/chat/send-generate-new-session", session.CreateGenerateSessionAndSendMessage)
	r.POST("/chat/send-generate", session.ChatGenerateSend)
	r.POST("/chat/send-stream-new-session", session.CreateStreamSessionAndSendMessage)
	r.POST("/chat/send-stream", session.ChatStreamSend)
	r.POST("/chat/history", session.ChatHistory)
	// r.POST("/chat/tts", AI.ChatSpeech)                  // ChatSpeechHandler
}
