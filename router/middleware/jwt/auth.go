package jwt

import (
	"GopherAI/common/code"
	response "GopherAI/router/controller"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
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

	log.Println("token is :", token)
	username, err := ParseToken(token)
	if !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
		// 验证不通过，直接断开
		c.Abort()
		return
	}

	c.Set("userName", username)
	c.Next()
}
*/

func Auth(c *gin.Context) {
	res := new(response.Response)

	// 从cookie中获取 access_token
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeNotLogin))
		c.Abort()
		return
	}

	log.Println("access token is :", token)

	deviceId := c.ClientIP() + c.GetHeader("user-agent")

	claims, err := ParseToken(token, deviceId)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidAccessToken))
		// 验证不通过，直接断开
		c.Abort()
		return
	}

	c.Set("userName", claims.Username)
	c.Next()
}
