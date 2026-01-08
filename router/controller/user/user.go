package user

import (
	"GopherAI/common/code"
	"GopherAI/router/controller"
	"GopherAI/router/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// =============login===============
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// omitempty当字段为空的时候，不返回这个东西
type LoginResponse struct {
	response.Response
	Token string `json:"token,omitempty"`
}

func Login(c *gin.Context) {
	req := new(LoginRequest)
	res := new(LoginResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Login(req.Username, req.Password)

	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// =============register===============
// 验证码由后端生成，存放到redis中，固然需要先发送一次请求CaptchaRequest,然后用返回的验证码
// 邮箱以及密码进行注册，后续再将账号进行返回
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

// 注册成功之后，直接让其进行登录状态
type RegisterResponse struct {
	response.Response
	Token string `json:"token,omitempty"`
}

func Register(c *gin.Context) {
	req := new(RegisterRequest)
	res := new(RegisterResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := user.Register(req.Email, req.Password, req.Captcha)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// =============captcha===============
type CaptchaRequest struct {
	Email string `json:"email" binding:"required"`
}

type CaptchaResponse struct {
	response.Response
}

func HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	code_ := user.SendCaptcha(req.Email)

	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	c.JSON(http.StatusOK, res)
}
