package model

import (
	"../data"
	"fmt"
)

type user struct {
	dataDB
}

var User user

func init() {
	User = user{dataDB{nil, false, data.UserInfo{}, data.UserInfoCreateTbl}}
	DBAddInit2Default(&User.dataDB)
}

type ShowUserInfo struct {
	ID       string `json:"id"`
	Account  string `json:"account"`  //账号
	Password string `json:"password"` //密码

	Name string `json:"name"` //住户名字
	Tel  string `json:"tel"`  //电话
}

func (ptr *user) GetShowInfo() ShowUserInfo {
	return ShowUserInfo{}
}

func (ptr *user) Add(show ShowUserInfo) error {
	if !ptr.isInit {
		return fmt.Errorf("pool not init")
	}
	db, err := ptr.pool.NewClient()
	if err != nil {
		return fmt.Errorf("pool new client err : %s", err)
	}
	defer db.Close()
	com := data.UserInfo{Account: show.Account, Password: show.Password, Name: show.Name, Tel: show.Tel}
	//参数判断
	if show.Account == "" || show.Password == "" {
		return fmt.Errorf("用户账号/密码不能为空")
	}
	return com.Add(db)
}
func (ptr *user) Info(id string) (show ShowUserInfo, err error) {
	if !ptr.isInit {
		err = fmt.Errorf("pool not init")
		return
	}
	db, err := ptr.pool.NewClient()
	if err != nil {
		err = fmt.Errorf("pool new client err : %s", err)
		return
	}
	defer db.Close()
	obj := data.UserInfo{ID: id}
	err = obj.Get(db)
	if err != nil {
		return
	}
	show.ID = obj.ID
	show.Account = obj.Account
	show.Password = obj.Password
	show.Tel = obj.Tel
	show.Name = obj.Name

	return
}
func (ptr *user) Update(id string, show ShowUserInfo) error {
	if !ptr.isInit {
		return fmt.Errorf("pool not init")
	}
	db, err := ptr.pool.NewClient()
	if err != nil {
		return fmt.Errorf("pool new client err : %s", err)
	}
	defer db.Close()
	com := data.UserInfo{ID: id, Name: show.Name, Tel: show.Tel,Password:show.Password}

	return com.Update(db)
}
func (ptr *user) Login(account, password string) (id string, ok bool, err error) {
	if !ptr.isInit {
		err = fmt.Errorf("pool not init")
		return
	}
	db, err := ptr.pool.NewClient()
	if err != nil {
		err = fmt.Errorf("pool new client err : %s", err)
		return
	}
	defer db.Close()

	return data.UserInfo{Account: account, Password: password}.Login(db)
}
