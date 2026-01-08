package router

import (
	"GopherAI/router/controller/user"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup) {
	r.POST("/register", user.Register)
	r.POST("/login", user.Login)
	r.POST("/captcha", user.HandleCaptcha)
}
