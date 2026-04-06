package device

import (
	"GopherAI/common/code"
	"GopherAI/common/mysql/dao/device"
	"GopherAI/common/mysql/model"
	"GopherAI/common/redis"
	"GopherAI/router/middleware/jwt"
	"time"
)

// 传入username 以备设备管理需要
func Logout(username string, accessToken string, refreshToken string, deviceId string) (error, code.Code) {
	atExpire, err := jwt.GetTokenExpireTime(accessToken)
	if err != nil {
		return err, code.CodeServerBusy
	}
	err = redis.AddToBlackList(accessToken, atExpire.Sub(time.Now()))
	if err != nil {
		return err, code.CodeServerBusy
	}

	rtExpire, err := jwt.GetTokenExpireTime(refreshToken)
	if err != nil {
		return err, code.CodeServerBusy
	}
	err = redis.AddToBlackList(refreshToken, rtExpire.Sub(time.Now()))
	if err != nil {
		return err, code.CodeServerBusy
	}

	// 从数据库device中删除
	err = device.DeleteDevice(username, deviceId)
	if err != nil {
		return err, code.CodeServerBusy
	}

	return nil, code.CodeSuccess
}

func DeviceList(username string) ([]model.DeviceInfo, code.Code) {
	devices, err := device.GetAllDeviceByUsername(username)
	if err != nil {
		return nil, code.CodeServerBusy
	}

	var deviceInfos []model.DeviceInfo

	for _, d := range devices {
		deviceInfos = append(deviceInfos, model.DeviceInfo{
			Username:      d.Username,
			DeviceIp:      d.DeviceIp,
			DeviceBrowser: d.DeviceBrowser,
			LoginTime:     d.IssuedTime,
		})
	}

	return deviceInfos, code.CodeSuccess
}

func OfflineDevice(username string, deviceIp string, deviceBrowser string) code.Code {
	// 1.从数据库中拿到token
	accessToken, refreshToken, err := device.GetAtAndRt(username, deviceIp, deviceBrowser)
	if err != nil {
		return code.CodeServerBusy
	}
	// 2.从数据库中删除device
	err = device.DeleteDeviceByIpAndBrowser(username, deviceIp, deviceBrowser)
	if err != nil {
		return code.CodeServerBusy
	}

	// 加入黑名单
	atExpire, err := jwt.GetTokenExpireTime(accessToken)
	if err != nil {
		return code.CodeServerBusy
	}
	err = redis.AddToBlackList(accessToken, atExpire.Sub(time.Now()))
	if err != nil {
		return code.CodeServerBusy
	}

	rtExpire, err := jwt.GetTokenExpireTime(refreshToken)
	if err != nil {
		return code.CodeServerBusy
	}
	err = redis.AddToBlackList(refreshToken, rtExpire.Sub(time.Now()))
	if err != nil {
		return code.CodeServerBusy
	}
	return code.CodeSuccess
}
