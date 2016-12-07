package generalcaptcha

import (
	"gopkg.in/redis.v5"
	"strconv"
	"time"
)

type RedisStore struct {
	Options map[string]interface{}
	RedisKeyPrefix string
	client *redis.Client
}

func (redisStore *RedisStore) InitialClient() {
	var options = &redis.Options{}
	SetFields(options, redisStore.Options)
	redisStore.client = redis.NewClient(options)
}

func (redisStore *RedisStore) Store(mobile, captcha string, timeout time.Duration) interface{} {
	reply := redisStore.client.Set(redisStore.RedisKeyPrefix + ":" + mobile, captcha, timeout)
	return reply
}

func (redisStore *RedisStore) Get(mobile string) string {
	reply, err := redisStore.client.Get(redisStore.RedisKeyPrefix + ":" + mobile).Result()
	if err != nil { // handle command error
		reply = ""
	}
	return reply
}

func (redisStore *RedisStore) IncreCountInDay(mobile string) {
	_, err := redisStore.client.HIncrBy(redisStore.RedisKeyPrefix, mobile, 1).Result()
	if err != nil {
		// handle command error
	}
}

func (redisStore *RedisStore) GetCountInDay(mobile string) (replyInt int) {
	reply, err := redisStore.client.HGet(redisStore.RedisKeyPrefix, mobile).Result()
	if err != nil {
		replyInt = 0
	} else {
		replyInt, _ = strconv.Atoi(reply)
	}
	return
}
