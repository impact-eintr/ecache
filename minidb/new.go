package minidb

func NewCache(dir string, ttl int) *MiniDB {
	db, err := Open(dir)
	if err != nil {
		return nil
	}
	return db
}
