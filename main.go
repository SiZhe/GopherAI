package main

import (
	"GopherAI/common/mysql"
	"GopherAI/common/rabbitmq"
	"GopherAI/common/redis"
	"GopherAI/config"
	"GopherAI/router"
	"log"
)

//func main() {
//	config.GetConfig()
//	err := rag.CreateUploadsDir("1765477265", "1")
//	if err != nil {
//		return
//	}
//
//	//StoreUploadsFiles()
//	pwd, _ := os.Getwd()
//	fmt.Println("pwd:%s", pwd)
//	filePath := filepath.Join("uploads_files", "1765477265", "1", "document.md")
//	fmt.Println("filePath:%s", filePath)
//
//	docs, err := rag.TransformerUploadsFiles(context.Background(), filePath)
//	if err != nil {
//		return
//	}
//
//	err = rag.EmbeddingUploadsFiles(context.Background(), docs)
//	if err != nil {
//		return
//	}
//}

func main() {
	//初始化config和env
	config.GetConfig()
	//初始化mysql
	err := mysql.InitMysql()
	if err != nil {
		log.Println("InitMysql error , " + err.Error())
		return
	} else {
		log.Println("mysql init success  ")
	}

	//初始化redis
	redis.InitRedis()
	log.Println("redis init success  ")

	//初始化rabbitmq
	rabbitmq.InitRabbitMQ()
	log.Println("rabbitmq init success  ")

	//同步数据库并启动HTTP服务
	err = router.InitServer()
	if err != nil {
		log.Println("GinServer error , " + err.Error())
		return
	} else {
		log.Println("gin init success  ")
	}

	rabbitmq.DestroyRabbitMQ()
}
