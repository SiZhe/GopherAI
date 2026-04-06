package router

import (
	"GopherAI/router/controller/device"

	"github.com/gin-gonic/gin"
)

func DeviceRouter(r *gin.RouterGroup) {
	r.POST("/logout", device.Logout)
	r.POST("/device-list", device.DeviceList)
	r.POST("/offline-device", device.OfflineDevice)
}
