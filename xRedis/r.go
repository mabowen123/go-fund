package xRedis

import (
	"fund/config"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const Ttl = time.Minute * 60 * 24 * 30

var client *redis.Client

func init() {
	redisConfig := config.Load("redis")
	dbIndex, _ := strconv.Atoi(redisConfig["db"])
	client = redis.NewClient(&redis.Options{
		Network:  redisConfig["network"],
		Addr:     redisConfig["address"] + ":" + redisConfig["port"],
		Password: redisConfig["password"],
		DB:       dbIndex,
	})
}

func HmSet(key string, value map[string]interface{}, ttl time.Duration) {
	client.HMSet(key, value)
	client.Expire(key, ttl)
}
func HGetAll(key string) (res map[string]string) {
	data := client.HGetAll(key)
	res, _ = data.Result()
	return
}

func Set(key string, value interface{}, ttl time.Duration) {
	client.Set(key, value, ttl)
}

func Get(key string) (res string) {
	res, _ = client.Get(key).Result()
	return
}
