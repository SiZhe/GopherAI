package router

import (
	"GopherAI/common/aihelper"
	"GopherAI/common/mysql/dao/message"
	"GopherAI/config"
	"GopherAI/router/middleware/jwt"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// 从数据库加载消息并初始化 AIHelperManager
func readDataFromDB() error {
	// 从数据库读取所有消息
	msgs, err := message.GetAllMessages()
	if err != nil {
		return err
	}

	// 遍历数据库消息
	for i := range msgs {
		msg := &msgs[i]
		config_ := map[string]string{
			"username":  msg.UserName,
			"sessionId": msg.SessionID,
		}

		// 创建对应的 AIHelper
		helper, err := aihelper.GetGlobalManager().GetOrCreateAIHelper(msg.UserName, msg.SessionID, msg.ModelType, config_)
		if err != nil {
			log.Printf("[readDataFromDB] failed to create helper for user=%s session=%s: %v", msg.UserName, msg.SessionID, err)
			continue
		}
		log.Println("readDataFromDB init:  ", helper.SessionID)
		// 添加消息到内存中(不开启存储功能)
		helper.AddMessage(msg.Content, msg.UserName, msg.IsUser, msg.ModelType, false)
	}

	log.Println("AIHelperManager init success ")
	return nil
}

func initRouter() *gin.Engine {
	r := gin.Default()
	//gin.SetMode("release")

	enterRouter := r.Group("/api/v2")
	{
		UserGroup := enterRouter.Group("/user")
		UserRouter(UserGroup)
	}

	//后续登录的接口需要jwt鉴权
	{
		AIGroup := enterRouter.Group("/AI")
		AIGroup.Use(jwt.Auth)
		AIRouter(AIGroup)
	}

	{
		FileGroup := enterRouter.Group("/file")
		FileGroup.Use(jwt.Auth)
		FileRouter(FileGroup)
	}

	{
		DeviceGroup := enterRouter.Group("/device")
		DeviceGroup.Use(jwt.Auth)
		DeviceRouter(DeviceGroup)
	}

	return r
}

//{
//	ImageGroup := enterRouter.Group("/image")
//	ImageGroup.Use(jwt.Auth())
//	ImageRouter(ImageGroup)
//}

func InitServer() error {
	// 首先同步数据库
	err := readDataFromDB()
	if err != nil {
		return err
	}

	r := initRouter()

	return r.Run(fmt.Sprintf("%s:%d", config.GetConfig().MainConfig.Host, config.GetConfig().MainConfig.Port))
}
