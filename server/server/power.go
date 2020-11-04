package server

import (
	"fmt"
	"sync"
)

/******************************* 权限管理相关 **************************************/
type powerVal struct {
	Power bool   `json:"power"`
	Name  string `json:"name"`  //功能名称
	State int    `json:"state"` //权限值 0/1/2、 默认、关闭、开发、
}

var defaultPower = map[string]*powerVal{} //默认权限定义
var initDefaultPowerMutex sync.Mutex

//注册默认访问权限
//path:uri路径
//power权限,true:默认有权限可以访问,false:默认没有权限不允许访问
//name该权限名称
func registerDefaultPower(path, name string, power bool) {
	initDefaultPowerMutex.Lock()
	defer initDefaultPowerMutex.Unlock()

	if path == "" {
		panic("path is nil")
	}

	val, exist := defaultPower[path]
	if exist {
		panic(fmt.Sprintf("%v is exist, val is : %+v", path, val))
	}

	fmt.Printf("Power register : %v -> %v : %v\n", path, name, power)
	defaultPower[path] = &powerVal{Power: power, Name: name, State: 0}
}


//检查权限,用户自身的权限高于默认权限
//如果用户自身有某一个接口的权限定义那么就会使用该定义,不会使用默认定义
//如果用户自身没有权限,就查看默认权限
//如果无默认权限则可以直接访问
func CheckPower(path string) (ok bool, err error) {

	//检查默认权限
	power, exist := defaultPower[path]
	if exist {
		if !power.Power {
			fmt.Println("检查权限:默认无权限")
			ok = false
			return
		} else {
			ok = true
			return
		}
	}

	//用户和默认都没有指明权限
	ok = true
	return
}

func getDefaultPowerList() (list map[string]*powerVal) {
	list = defaultPower
	return
}