package model

import (
	"fmt"
	"github.com/seefan/gossdb"
	"github.com/seefan/gossdb/conf"
)

var pool *gossdb.Connectors = nil

var devPool *gossdb.Connectors = nil

//初始化数据库
func InitDB(host string ,port int,pwd string) (err error) {

	fmt.Println("init DB start... host/post",host,port)
	pool, err = gossdb.NewPool(&conf.Config{
		Host:             host,
		Port:             port,
		Password:pwd,
		RetryEnabled:     true, //是否重连
		MinPoolSize:      20,
		MaxPoolSize:      2000,
		MaxWaitSize:      10000,
		AcquireIncrement: 5,
	})
	if err != nil {
		fmt.Println("create pool error", err)
		return
	}

	return nil
}
func CloseDB() {
	if pool != nil {
		pool.Close()
	}
}
//初始化数据库
func InitBaseDeviceSSDB(host string ,port int,pwd string) (err error) {

	fmt.Println("init base device SSDB start...")
CONN :

	devPool, err = gossdb.NewPool(&conf.Config{
		Host:             host,
		Port:             port,
		Password:		  pwd,
		RetryEnabled:     true, //是否重连
		MinPoolSize:      20,
		MaxPoolSize:      100,
		MaxWaitSize:      10000,
		AcquireIncrement: 5,
	})

	if err != nil {
		fmt.Printf("create pool error %s re connct \n", err)
		goto  CONN
	}
	return nil
}
func CloseBaseDeviceDB() {
	if devPool != nil {
		devPool.Close()
	}
}