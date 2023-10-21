package main

import (
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

func getSecondOfDayBetween() time.Duration {
	t := time.Now()
	return time.Duration(float64(60*60*24-60*60*t.Hour()-60*t.Minute()-t.Second()) * float64(time.Second))
}

func SetCachedValue(url string, dataJson string) {
	err := rdb.Set(ctx, url, dataJson, getSecondOfDayBetween()).Err()
	if err != nil {
		panic(err)
	}
}

func GetCachedValue(url string) (string, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := rdb.Get(ctx, url).Result()
	if err != redis.Nil && err != nil {
		panic(err)
	}

	return val, err
}
