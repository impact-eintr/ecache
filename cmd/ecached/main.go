package main

import (
	"flag"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/impact-eintr/ecache"
	"github.com/impact-eintr/ecache/http"
	"github.com/impact-eintr/ecache/tcp"
)

func main() {

	var typ, cacheDir string
	var httpPort, tcpPort string
	var ttl int

	flag.StringVar(&typ, "t", "mem", "缓存类型 可选项 mem disk")

	flag.StringVar(&tcpPort, "tp", "6430", "ecached tcp 服务端口")
	flag.StringVar(&httpPort, "hp", "7895", "ecached http 服务端口")

	flag.StringVar(&cacheDir, "d",
		func() string {
			pwd, _ := os.Getwd()
			return path.Join(pwd, "storage")
		}(), "磁盘缓存目录")
	flag.IntVar(&ttl, "T", 0, "缓存生存时间 默认为0 即不失效 采用LRU算法")

	flag.Parse()

	c := ecache.New(typ)
	go tcp.New(c).Listen()

	router := gin.Default()

	ch := http.NewCacheHandler(c)
	router.GET("/cache/*key", ch.GetHandler)
	router.PUT("/cache/*key", ch.PutHandler)
	router.DELETE("/cache/*key", ch.DelHandler)

	router.GET("/status", ch.StatusHandler)

	router.Run(":" + httpPort)

}