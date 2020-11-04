package model

import (
	"../data"
	"github.com/seefan/gossdb"
	"sync"
)

type dataDB struct {
	pool      *gossdb.Connectors
	isInit    bool
	table     data.Table
	createTbl func(pool *gossdb.Connectors)
}

var defaultPool *gossdb.Connectors = nil
var defaultDBInitList []*dataDB
var defaultDBMutex sync.Mutex

func DBAddInit2Default(dataDb *dataDB) {
	defaultDBMutex.Lock()
	defer defaultDBMutex.Unlock()
	defaultDBInitList = append(defaultDBInitList, dataDb)
}