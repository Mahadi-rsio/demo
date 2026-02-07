package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func Redis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping().Result()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(pong)

	return rdb

}
