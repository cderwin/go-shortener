package main

import (
	"github.com/gocraft/web"
	"github.com/patrickmn/go-cache"
	"time"
)

func NewMockServer() Server {
	mockRedis, _ := CreateMockStore()
	cache := cache.New(5*time.Minute, 30*time.Second)
	return Server{cache, mockRedis}
}

func NewMockRouter() *web.Router {
	server := NewMockServer()
	router := web.New(server)
	setupRoutes(router, server)
	return router
}
