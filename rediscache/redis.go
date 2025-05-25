package rediscache

import (
	"context"
	"stock-analyser/logger"
	"stock-analyser/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var Client *redis.Client

func Redis() {
	logger.Infof("Initializing redis connection...: RedisServerHost=%s, RedisServerPort=%s", utils.Config.RedisServerHost, utils.Config.RedisServerPort)
	opt, _ := redis.ParseURL("rediss://" + utils.Config.RedisServerHost + ":" + utils.Config.RedisServerPort)
	Client = redis.NewClient(opt)

}

func AddDataToCache(username string, email string) {
	Client.Set(ctx, "user:"+username, "email:"+email, 20*time.Minute)
	Client.Set(ctx, "email:"+email, "user:"+username, 20*time.Minute)
}
