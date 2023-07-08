package store

import (
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

type LevelDBProvider struct {
	dbPath   string
	DB       *leveldb.DB
	lock     *sync.RWMutex // 读写锁
	NotFound error
}

func NewLevelDBProvider(filepath string) *LevelDBProvider {
	db, err := leveldb.OpenFile(filepath, nil)
	if err != nil {
		log.Errorf("leveldb open fiaild: %v", err)
		panic(err)
	}
	return &LevelDBProvider{
		dbPath:   filepath,
		DB:       db,
		lock:     new(sync.RWMutex),
		NotFound: leveldb.ErrNotFound,
	}
}

func (lp *LevelDBProvider) GetOne(key []byte) ([]byte, error) {
	lp.lock.RLock()
	defer lp.lock.RUnlock()
	return lp.DB.Get(key, nil)
}

func (lp *LevelDBProvider) Set(key []byte, value []byte) error {
	lp.lock.Lock()
	defer lp.lock.Unlock()
	return lp.DB.Put(key, value, &opt.WriteOptions{Sync: true})
}

func (lp *LevelDBProvider) Del(key []byte) (bool, error) {
	lp.lock.Lock()
	defer lp.lock.Unlock()
	err := lp.DB.Delete(key, &opt.WriteOptions{
		Sync: false,
	})
	res := true
	if err != nil {
		res = false
	}
	return res, err
}

func (lp *LevelDBProvider) Has(key []byte) (bool, error) {
	lp.lock.RLock()
	defer lp.lock.RUnlock()
	return lp.DB.Has(key, nil)
}
func (lp *LevelDBProvider) GetKeys() ([][]byte, error) {
	res := make([][]byte, 0)
	lp.lock.RLock()
	defer lp.lock.RUnlock()
	iter := lp.DB.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		res = append(res, iter.Key())
	}
	return res, nil
}

func (lp *LevelDBProvider) GetValues() ([][]byte, error) {
	res := make([][]byte, 0)
	lp.lock.RLock()
	defer lp.lock.RUnlock()
	iter := lp.DB.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		res = append(res, iter.Value())
	}
	return res, nil
}
