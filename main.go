package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocraft/web"
	"github.com/patrickmn/go-cache"
)

// Derivative router type to enable bulk addition of middleware

func main() {
	server := createServer()
	router := web.New(server)
	setupRoutes(router, server)
	router.Middleware(web.LoggerMiddleware)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRoutes(router *web.Router, server Server) {
	router.Get("/healthcheck", server.healthcheck)
	router.Post("/create", server.addUrl)
	router.Get("/:path", server.fetchUrl)
	router.Get("/stats/:path", server.urlStats)
}

type Server struct {
	UrlCache *cache.Cache
	Redis     Datastore
}

func createServer() Server {
	redisUrl := os.Getenv("REDIS_URL")
	redisClient := NewRedisStore(redisUrl)
	urlCache := cache.New(5*time.Minute, 30*time.Second)
	server := Server{UrlCache: urlCache, Redis: redisClient}
	return server
}
