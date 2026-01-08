package jwt

import (
	"GopherAI/common/code"
	"GopherAI/router/controller"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	res := new(response.Response)

	var token string

	authHeader := c.GetHeader("Authorization")

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		token = c.Query("token")
	}

	if token == "" {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
		// 验证不通过，直接断开
		c.Abort()
		return
	}

	log.Println("token is ", token)
	username, ok := ParseToken(token)
	if !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
		// 验证不通过，直接断开
		c.Abort()
		return
	}

	c.Set("userName", username)
	c.Next()
}
