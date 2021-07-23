package model

import (
	"day4/internal"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

// ManRedisKey redis接口
type ManRedisKey interface {
	StrUpdate(value interface{}) string //用户获取礼品并更新redis
	StrHave(value interface{}) bool     //判断数据库中是否有此条数据
}

//redis参数
const (
	RedisURL            = "redis://127.0.0.1:6379"
	redisMaxIdle        = 25  //最大空闲连接数
	redisMaxActive      = 100 //最大的激活连接数
	redisIdleTimeoutSec = 240 //最大空闲连接时间
	//RedisPassword       = ""
)

//创建用户
var RedisPool *redis.Pool
var RedisPool1 *redis.Pool

// NewRedisPool 新建redis池
func NewRedisPool(redisURL string, Database int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   redisMaxActive,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			// redis://127.0.0.1:6379/0
			URLs := strings.Join([]string{redisURL, strconv.Itoa(Database)}, "/")
			fmt.Println(URLs)
			c, err := redis.DialURL(URLs)
			//c, err := redis.Dial("tcp",redisURL,redis.DialDatabase(5))
			if err != nil {
				return nil, internal.InternalServiceError(err.Error())
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return internal.InternalServiceError("ping redis error: " + err.Error())
			}
			return nil
		},
	}
}
