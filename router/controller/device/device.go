package device

import (
	"GopherAI/common/code"
	"GopherAI/common/mysql/model"
	response "GopherAI/router/controller"
	"GopherAI/router/service/device"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// =============logout===============
type LogoutRequest struct {
	Username string `json:"username"`
}

type LogoutResponse struct {
	response.Response
}

func Logout(c *gin.Context) {
	req := new(LogoutRequest)
	res := new(LogoutResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	at, _ := c.Cookie("access_token")
	rt, _ := c.Cookie("refresh_token")

	userName := c.GetString("userName")

	deviceId := c.ClientIP() + c.GetHeader("user-agent")

	//fmt.Println("Username:", userName)

	err, code_ := device.Logout(userName, at, rt, deviceId)
	if err != nil {
		log.Println(err)
	}
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	// 清空Cookie
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	res.Success()
	c.JSON(http.StatusOK, res)
}

// =============deviceList===============
type DeviceListRequest struct {
	Username string `json:"username"`
}

type DeviceListResponse struct {
	Devices []model.DeviceInfo `json:"sessions,omitempty"`
	response.Response
}

func DeviceList(c *gin.Context) {
	req := new(DeviceListRequest)
	res := new(DeviceListResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	userName := c.GetString("userName")

	devices_, code_ := device.DeviceList(userName)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Devices = devices_
	res.Success()
	c.JSON(http.StatusOK, res)
}

// =============offlineDevice===============
type OfflineDeviceRequest struct {
	Username      string `json:"username"`
	DeviceIp      string `json:"device_ip"`
	DeviceBrowser string `json:"device_browser"`
}

type OfflineDeviceResponse struct {
	response.Response
}

func OfflineDevice(c *gin.Context) {
	req := new(OfflineDeviceRequest)
	res := new(OfflineDeviceResponse)

	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	username := c.GetString("userName")
	code_ := device.OfflineDevice(username, req.DeviceIp, req.DeviceBrowser)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	c.JSON(http.StatusOK, res)
}
