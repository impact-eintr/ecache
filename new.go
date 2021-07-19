package ecache

import (
	"log"

	"github.com/impact-eintr/ecache/cache"
	"github.com/impact-eintr/ecache/global"
	"github.com/impact-eintr/ecache/mem"
	"github.com/impact-eintr/ecache/rdb"
)

func New(typ string) (c cache.Cache) {

	if typ == "mem" {
		c = mem.NewCache(global.TTL)
	} else if typ == "disk" {
		c = rdb.NewCache(global.CacheDir, global.TTL)
	}

	if c == nil {
		panic("未指定类型")
	}

	log.Println(typ, "服务已就位")
	return c
}
