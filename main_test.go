package main

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func CreateMockServer() Server {
	mockRedis, _ := CreateMockStore()
	cache := cache.New(5*time.Minute, 30*time.Second)
	return Server{cache, mockRedis}
}
