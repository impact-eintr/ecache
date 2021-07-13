package cache

type Stat struct {
	Count     int64 `json:"count"`
	KeySize   int64 `json:"keysize"`
	ValueSize int64 `json:"valuesize"`
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))

}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
