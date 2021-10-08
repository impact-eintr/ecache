package minidb

import (
	"io"
	"os"
	"sync"

	"github.com/impact-eintr/ecache/cache"
)

type MiniDB struct {
	indexes map[string]int64 // 内存中的索引信息
	dbFile  *DBFile          // 数据文件
	dirPath string
	mu      sync.RWMutex
}

// Open 开启一个数据库实例
func Open(dirPath string) (*MiniDB, error) {
	// 如果数据库目录不存在 新建
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 加载数据文件
	dbFile, err := NewDBFile(dirPath)
	if err != nil {
		return nil, err
	}

	db := &MiniDB{
		dbFile:  dbFile,
		indexes: make(map[string]int64),
		dirPath: dirPath,
	}

	// 加载索引
	db.loadIndexesFromFile(dbFile)
	return db, nil
}

func (db *MiniDB) Set(key string, value []byte) (err error) {
	if len(key) == 0 {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	offset := db.dbFile.Offset
	// 封装成 Entry
	entry := NewEntry([]byte(key), value, PUT)
	// 追加到数据文件中
	err = db.dbFile.Write(entry)
	if err != nil {
		return
	}

	// 写到内存
	db.indexes[key] = offset
	return

}

func (db *MiniDB) Get(key string) (val []byte, err error) {
	if len(key) == 0 {
		return
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	// 从内存中取出索引
	offset, ok := db.indexes[key]
	if !ok {
		return
	}

	// 从磁盘中读取数据
	var e *Entry
	e, err = db.dbFile.Read(offset)
	if err != nil && err != io.EOF {
		return
	}
	if e != nil {
		val = e.Value
	}
	return
}

// Del 删除数据
// 这里并不会定位到原记录进行删除，而还是将删除的操作封装成 Entry，
// 追加到磁盘文件当中，只是这里需要标识一下 Entry 的类型是删除。
func (db *MiniDB) Del(key string) (err error) {
	if len(key) == 0 {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.indexes[key]
	if !ok {
		return
	}

	// 封装成 Entry
	entry := NewEntry([]byte(key), nil, DEL)
	// 追加到数据文件中
	err = db.dbFile.Write(entry)
	if err != nil {
		return
	}

	// 从内存删除
	delete(db.indexes, string(key))
	return

}

func (db *MiniDB) GetStat() cache.Stat {
	return cache.Stat{}
}

// Merge合并数据文件
func (db *MiniDB) Merge() error {
	// 没有数据
	if db.dbFile.Offset == 0 {
		return nil
	}

	var (
		validEntries []*Entry
		offset       int64
	)

	// 读取元数据文件的Entry
	for {
		e, err := db.dbFile.Read(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// 内存中的索引是最新的，直接对比过滤出有效的Entry
		if off, ok := db.indexes[string(e.Key)]; ok && off == offset {
			validEntries = append(validEntries, e)
		}
		offset += e.GetSize()
	}
	if len(validEntries) > 0 {
		// 新建临时文件
		mergeDBFile, err := NewMergeDBFile(db.dirPath)
		if err != nil {
			return err
		}
		defer os.Remove(mergeDBFile.File.Name())

		// 重新写入有效的 entry
		for _, entry := range validEntries {
			writeOff := mergeDBFile.Offset
			err := mergeDBFile.Write(entry)
			if err != nil {
				return err
			}

			// 更新索引
			db.indexes[string(entry.Key)] = writeOff
		}
		// 获取文件名
		dbFileName := db.dbFile.File.Name()
		// 关闭文件
		db.dbFile.File.Close()
		// 删除旧的数据文件
		os.Remove(dbFileName)

		// 获取文件名
		mergeDBFileName := mergeDBFile.File.Name()
		// 关闭文件
		mergeDBFile.File.Close()
		// 临时文件变更为新的数据文件
		os.Rename(mergeDBFileName, db.dirPath+string(os.PathSeparator)+FileName)

		db.dbFile = mergeDBFile

	}
	return nil
}

func (db *MiniDB) loadIndexesFromFile(dbFile *DBFile) {
	if dbFile == nil {
		return
	}

	var offset int64
	for {
		// offset 传递的是值
		e, err := db.dbFile.Read(offset)
		if err != nil {
			// 读取完毕
			if err == io.EOF {
				break
			}
			return
		}

		// 设置索引状态 这里其实可以尝试使用bigcache的设计
		db.indexes[string(e.Key)] = offset

		if e.Mark == DEL {
			// 删除内存中的 key
			delete(db.indexes, string(e.Key))
		}
		offset += e.GetSize()
	}
	return
}
