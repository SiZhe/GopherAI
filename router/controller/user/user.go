package user

import (
	"GopherAI/common/code"
	"GopherAI/config"
	"GopherAI/router/controller"
	"GopherAI/router/service/user"
	"net/http"
	"time"

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
	//AccessToken  string `json:"accessToken,omitempty"`
	//RefreshToken string `json:"refreshToken,omitempty"`
}

func Login(c *gin.Context) {
	req := new(LoginRequest)
	res := new(LoginResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	deviceId := c.ClientIP() + c.GetHeader("user-agent")

	accessToken, refreshToken, code_ := user.Login(req.Username, req.Password, deviceId)

	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute
	refreshExpireDuration := time.Duration(config.GetConfig().JwtConfig.RefreshExpire) * time.Hour

	// 设置Cookie（本地 Secure=false ，上线改为true）
	// 参数: 名称，值，时间，路径 ，域名，secure: 是否只走 https ,httponly：前端无法获取
	c.SetCookie("access_token", accessToken, int(accessExpireDuration.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, int(refreshExpireDuration.Seconds()), "/", "", false, true)

	res.Success()
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
	// Token string `json:"token,omitempty"`
}

func Register(c *gin.Context) {
	req := new(RegisterRequest)
	res := new(RegisterResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	deviceId := c.ClientIP() + c.GetHeader("user-agent")

	accessToken, refreshToken, code_ := user.Register(req.Email, req.Password, req.Captcha, deviceId)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute
	refreshExpireDuration := time.Duration(config.GetConfig().JwtConfig.RefreshExpire) * time.Hour

	// 设置Cookie（本地 Secure=false ，上线改为true）
	// 参数: 名称，值，时间，路径 ，域名，secure: 是否只走 https ,httponly：前端无法获取
	c.SetCookie("access_token", accessToken, int(accessExpireDuration.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, int(refreshExpireDuration.Seconds()), "/", "", false, true)

	res.Success()
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

// =============reflashAT===============
type ReflashATRequest struct {
	Username string `json:"username"`
}

type ReflashATResponse struct {
	response.Response
}

func RefreshAccessToken(c *gin.Context) {
	req := new(ReflashATRequest)
	res := new(ReflashATResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidRefreshToken))
		return
	}

	deviceId := c.ClientIP() + c.GetHeader("user-agent")
	accessExpireDuration := time.Duration(config.GetConfig().JwtConfig.AccessExpire) * time.Minute

	newAccessToken, code_ := user.RefreshAccessTokenWithToken(rt, deviceId)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	c.SetCookie("access_token", newAccessToken, int(accessExpireDuration.Seconds()), "/", "", false, true)

	res.Success()
	c.JSON(http.StatusOK, res)
}
