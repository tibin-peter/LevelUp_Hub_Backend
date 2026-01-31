package redis

import (
	"LevelUp_Hub_Backend/internal/config"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(cfg *config.Config)(*redis.Client,error){
	rdb:=redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		ReadTimeout: 5 *time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize: 10,//for pooling
	})
	_,err:=rdb.Ping(Ctx).Result()//text command
	if err!=nil{
		return nil,err
	}

	log.Println("Connected to redis")

	return rdb,nil
}