package cache

import (
	"log"
)

func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "mem" {
		c = newMemCache(ttl)
	}

	if c == nil {
		log.Fatalln("unknown cache type " + typ)
	}
	return c
}
