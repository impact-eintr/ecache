package cache

import "log"

type Stat struct {
	Count     int64 `json:"count"`
	KeySize   int64 `json:"keysize"`
	ValueSize int64 `json:"valuesize"`
}

func (s *Stat) Add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
	log.Println("add()", s)
}

func (s *Stat) Remove(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
	log.Println("remove()", s)
}
