package data

import (
	"fmt"
	"github.com/seefan/gossdb"
	"github.com/seefan/sorm"
)

type Users1 struct {
	UsersBean1
	Job  string `json:"job"`
	Sex  bool 	`json:"sex"`
	Avg  int64  `json:"avg,omitempty"`
	Addr string `json:"addr,omitempty"`
}
type UsersBean1 struct {
	Id   string	`json:"id",tbl:"KEY"`
	Name string	`json:"name"`
	Age  int    `json:"age"`

}
/* 没有tbl tag 默认用DATA */
type users1Cols struct {
	Name string `tbl:"!DATA"`
	Age  string `tbl:"!NUM"`
	Sex  string `tbl:"!NUM"`
	Job  string `tbl:"DATA"`
	Avg  string `tbl:"NUM"`
	Addr string `tbl:"!DATA"`
}

//自定义类的各字段对应的数据库表的列名，优先级高于类json列名设置，当都没有设置是，默认为字段名
var users1Colsv = &users1Cols{
	/*
		Name: "NAME",
		Age:  "age",
		Job:  "job",
		Avg:  "avg",
		Addr: "addr",
	*/
}

//指定排序字段
type users1Sort struct {
	Age string
	Avg string
}

var user1Sortv = &users1Sort{}

/* 如果没有指定 cidx 默认是非唯一  */
type users1Cidx struct {
	Name              []string `cidx:"only"`
	Addr              []string
	Age               []string
	Sex               []string
	AddrName          []string
	JobAddrNameAvg    []string
	JobAddrNameAvgAge []string
	Test              []string
}

//自定义索引
var users1Cidxv = &users1Cidx{
	AddrName:   []string{"addr", "name"},
	Test: 		[]string{"name", "addr"},
}

//自定义数据库表名
func (u *Users1) TblName() string {
	return "users12334"
}

func (u *Users1) TblType() string {
	return "INCR"
}

func (u *Users1) Cols() *users1Cols {
	return users1Colsv
}

func (u *Users1) Sort() *users1Sort {
	return user1Sortv
}

func (u *Users1) Cidx() *users1Cidx {
	return users1Cidxv
}

func (u *Users1) Print2Map() (uinfo map[string]interface{} ){

	uCol := u.Cols()

	uinfo = make(map[string]interface{})

	uinfo[uCol.Addr] = u.Addr
	uinfo[uCol.Name] = u.Name
	uinfo[uCol.Age] = u.Age
	uinfo[uCol.Avg] = u.Avg
	uinfo[uCol.Job] = u.Job

	return
}

func (u *Users1)Add(hashName string,db *gossdb.Client) (id string,err error){

	orm := sorm.Engine{DB:db}

	orm.SetTable(sorm.Tables[u.TblName()])

	ids,err := orm.Htinsert(hashName,u)

	if err != nil{
		return "", fmt.Errorf(" Add error: %s",err.Error())
	}
	if len(ids) != 1{
		return "", fmt.Errorf(" Add error")
	}
	return ids[0],nil
}
func (u *Users1)Del(hashName string,db *gossdb.Client) (ok bool,err error){
	if u.Id != ""{
		ret,err := db.Ht_row_clear(u.TblName(),hashName,u.Id)
		if err != nil {
			return false, fmt.Errorf(" Del error: %s",err.Error())
		}
		if ret == 0 {
			return true,fmt.Errorf(" %s not exists",u.Id)
		}
		return true,nil
	}else {
		_,_,_,err = db.Ht_clear(u.TblName(),hashName)

		return err == nil,err
	}
}
//当id空是，表示查询所有信息
//当id不为空，表示查询指定信息
func (u *Users1)Get(hashName string,db *gossdb.Client,cols ...string)( us []Users1,err error)  {

	if u.Id == ""{
		ret,err := db.Ht_row_scan(u.TblName(),hashName,"","",-1,"*")

		if err != nil {
			return nil, fmt.Errorf(" Get error:%s",err.Error())
		}
		us = make([]Users1, len(ret))
			for k,v := range ret{
				node := Users1{}

				node.Id = v.Id
				node.Age = v.Colv[node.Cols().Age].Int()
				node.Sex = v.Colv[node.Cols().Sex].Bool()
				node.Name = v.Colv[node.Cols().Name].String()
				node.Addr = v.Colv[node.Cols().Addr].String()
				node.Job = v.Colv[node.Cols().Job].String()
				node.Avg = v.Colv[node.Cols().Avg].Int64()

				us[k] = node
			}
			*u = us[0]
	}else{
		if len(cols) == 1 && cols[0] == "*"{
			cols = make([]string,0)
			cols = append(cols, u.Cols().Avg)
			cols = append(cols, u.Cols().Job)
			cols = append(cols, u.Cols().Addr)
			cols = append(cols, u.Cols().Sex)
			cols = append(cols, u.Cols().Age)
			cols = append(cols, u.Cols().Name)
		}

		//Multi_ht_col_get 不支持 * 必须指定哪些列名
		ret,err := db.Multi_ht_col_get(u.TblName(),hashName,u.Id,cols...)

		if err != nil {
			return nil, fmt.Errorf(" Get error:%s",err.Error())
		}
		if len(ret) != len(cols){
			return nil, fmt.Errorf(" Get info")
		}
		for k,v := range cols{
			switch v {
			case u.Cols().Avg:
				u.Avg = ret[k].Int64()
			case u.Cols().Job:
				u.Job = ret[k].String()
			case u.Cols().Addr:
				u.Addr = ret[k].String()
			case u.Cols().Sex:
				u.Sex = ret[k].Bool()
			case u.Cols().Age:
				u.Age = ret[k].Int()
			case u.Cols().Name:
				u.Name = ret[k].String()
			}
		}
		us = append(us, *u)
	}
	return
}

func init()  {

	err := sorm.Sync(&Users1{})
	if err != nil {
		fmt.Printf("sorm sync err:%s\n",err.Error())
	}

}