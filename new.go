package ecache

import (
	"log"

	"github.com/impact-eintr/ecache/cache"
	"github.com/impact-eintr/ecache/mem"
	"github.com/impact-eintr/ecache/rdb"
)

var (
	CacheDir string = ""
	TTL      int    = 0
	TcpPort  string = "6430"
)

func New(typ string) (c cache.Cache) {

	if typ == "mem" {
		c = mem.NewCache(TTL)
	} else if typ == "disk" {
		c = rdb.NewCache(CacheDir, TTL)
	}

	if c == nil {
		panic("未指定类型")
	}

	log.Println(typ, "服务已就位")
	return c
}
