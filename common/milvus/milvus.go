package milvus

import (
	"GopherAI/config"
	"context"
	"sync"

	milvusCli "github.com/milvus-io/milvus-sdk-go/v2/client"
)

var milvusClient milvusCli.Client

var once sync.Once

func GetMilvusClient() milvusCli.Client {
	once.Do(func() {
		ctx := context.Background()
		cli, err := milvusCli.NewClient(ctx, milvusCli.Config{
			Address: config.GetConfig().MilvusConfig.MilvusAddress,
			DBName:  config.GetConfig().MilvusConfig.MilvusDb,
		})
		if err != nil {
			panic(err)
		}
		milvusClient = cli
	})
	return milvusClient
}
