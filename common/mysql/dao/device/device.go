package device

import (
	"GopherAI/common/mysql"
	"GopherAI/common/mysql/dao/user"
	"GopherAI/common/mysql/model"
	"GopherAI/utils"
	"strings"
	"time"
)

func CreateDevice(username string, deviceInfo string, issuedTime time.Time, at string, rt string) (*model.Device, error) {
	if strings.Contains(username, "@") {
		// 传入的是邮箱
		err, un := user.GetUsernameByEmail(username)
		if err != nil {
			return nil, err
		}
		username = un
	}

	ip, browser := utils.ParseIPAndBrowser(deviceInfo)

	device := &model.Device{
		Username:      username,
		DeviceIp:      ip,
		DeviceBrowser: browser,
		DeviceInfo:    deviceInfo,
		AccessToken:   at,
		RefreshToken:  rt,
		IssuedTime:    issuedTime,
	}

	err := mysql.DB.Create(&device).Error

	if err != nil {
		return nil, err
	}

	return device, nil
}

func GetDevicesByUser(username string) ([]model.Device, error) {
	var devices []model.Device

	// 按用户名查询，按签发时间倒序（最新的在前）
	err := mysql.DB.Where("username = ?", username).
		Order("issued_time DESC").
		Find(&devices).Error

	if err != nil {
		return nil, err
	}
	return devices, nil
}

func DeleteDevice(username string, deviceInfo string) error {
	if strings.Contains(username, "@") {
		// 传入的是邮箱
		err, un := user.GetUsernameByEmail(username)
		if err != nil {
			return err
		}
		username = un
	}
	//fmt.Println("username:", username)
	//fmt.Println("deviceInfo:", deviceInfo)

	ip, browser := utils.ParseIPAndBrowser(deviceInfo)

	err := mysql.DB.Where("username = ? AND device_ip = ? AND device_browser = ?", username, ip, browser).
		Unscoped(). // 关键：强制物理删除
		Delete(&model.Device{}).Error

	return err
}

func DeleteDeviceByIpAndBrowser(username string, ip string, browser string) error {
	if strings.Contains(username, "@") {
		// 传入的是邮箱
		err, un := user.GetUsernameByEmail(username)
		if err != nil {
			return err
		}
		username = un
	}

	err := mysql.DB.Where("username = ? AND device_ip = ? AND device_browser = ?", username, ip, browser).
		Unscoped(). // 关键：强制物理删除
		Delete(&model.Device{}).Error

	return err
}

func GetAllDeviceByUsername(username string) ([]model.Device, error) {
	// 如果是邮箱，转换成用户名
	if strings.Contains(username, "@") {
		err, un := user.GetUsernameByEmail(username)
		if err != nil {
			return nil, err
		}
		username = un
	}

	// 定义接收结果的切片
	var devices []model.Device

	// GORM 查询：根据用户名查找设备
	err := mysql.DB.Where("username = ?", username).Find(&devices).Error

	// 返回结果和错误
	return devices, err
}

func GetAtAndRt(username string, ip string, browser string) (string, string, error) {
	if strings.Contains(username, "@") {
		err, un := user.GetUsernameByEmail(username)
		if err != nil {
			return "", "", err
		}
		username = un
	}

	// 定义接收结果的切片
	var device model.Device

	err := mysql.DB.Where("username = ? AND device_ip = ? AND device_browser = ?", username, ip, browser).First(&device).Error

	return device.AccessToken, device.RefreshToken, err
}
