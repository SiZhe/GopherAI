package router

import (
	"GopherAI/router/controller/user"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	// 不需要认证的接口
	r.POST("/register", user.Register)
	r.POST("/login", user.Login)
	r.POST("/captcha", user.HandleCaptcha)
	r.POST("/refresh-access-token", user.RefreshAccessToken)
}
