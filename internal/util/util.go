package util

import model2 "day4/model"

func Init() {
	model2.RedisPool = model2.NewRedisPool(model2.RedisURL, 1)
	model2.Client = model2.GetClient()
}
