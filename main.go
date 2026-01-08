package main

import (
	"GopherAI/common/mysql"
	"GopherAI/common/rabbitmq"
	"GopherAI/common/redis"
	"GopherAI/router"
	"log"
)

func main() {
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
