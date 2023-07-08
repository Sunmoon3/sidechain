package store

import (
	log "github.com/sirupsen/logrus"
)

type StoreProvider interface {
	GetOne(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	Del(key []byte) (bool, error)
	Has(key []byte) (bool, error)
	GetKeys() ([][]byte, error)
	GetValues() ([][]byte, error)
}

type StoreService struct {
	storeProvider StoreProvider
}

func NewStoreService(sp StoreProvider) *StoreService {
	return &StoreService{storeProvider: sp}
}

func (ss *StoreService) QueryByKey(key []byte) []byte {
	val, err := ss.storeProvider.GetOne(key)
	if err != nil {
		log.Errorf("QueyByKey error: %v\n", err)
		return nil
	}
	return val
}

func (ss *StoreService) StoreKeyValue(key, value []byte) bool {
	err := ss.storeProvider.Set(key, value)
	if err != nil {
		log.Errorf("StoreKeyValue error: %v\n", err)
		return false
	}
	return true
}

func (ss *StoreService) DeleteByKey(key []byte) bool {
	res, err := ss.storeProvider.Del(key)
	if err != nil {
		log.Errorf("DeleteByKey error: %v\n", err)
		return false
	}
	return res
}

func (ss *StoreService) KeyExists(key []byte) bool {
	res, err := ss.storeProvider.Has(key)
	if err != nil {
		log.Errorf("KeyExists error: %v\n", err)
		return false
	}
	return res
}

func (ss *StoreService) GetAllKeyValues() ([][]byte, [][]byte) {
	keys, err := ss.storeProvider.GetKeys()
	if err != nil {
		log.Errorf("KeyExists error: %v\n", err)
		return nil, nil
	}
	values, err := ss.storeProvider.GetValues()
	if err != nil {
		log.Errorf("KeyExists error: %v\n", err)
		return nil, nil
	}
	return keys, values
}
