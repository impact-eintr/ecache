package main

import (
	"github.com/gin-gonic/gin"
	"github.com/impact-eintr/ecache/cache"
	"github.com/impact-eintr/ecache/http"
	"github.com/impact-eintr/ecache/tcp"
)

func main() {
	c := cache.New("mem", 5)
	go tcp.New(c).Listen()
	ch := http.NewCacheHandler(c)

	router := gin.Default()
	router.GET("/cache/:key", ch.GetHandler)
	router.PUT("/cache/:key", ch.PutHandler)
	router.DELETE("/cache/:key", ch.DelHandler)

	router.GET("/status", ch.StatusHandler)

	router.Run(":8080")

}
