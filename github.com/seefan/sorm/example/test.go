package main

import (
	"fmt"
	"github.com/seefan/gossdb/conf"
	"strings"

	"github.com/seefan/gossdb"
	"github.com/seefan/gossdb/ssdb"
)

/********
定义表结构
*********/
//json设置类字段对应在数据库中的列名

type users3 struct {
	Id   string	`json:"id",tbl:"key"`
	Name string	`json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
	Avg  int64  `json:"avg,omitempty"`
	Addr string `json:"omitempty"`
}

/* 没有tbl tag 默认用DATA */
type users3Cols struct {
	Name string `tbl:"!DATA"`
	Age  string `tbl:"!NUM"`
	Job  string `tbl:"!DATA"`
	Avg  string `tbl:"NUM"`
	Addr string
}

//自定义类的各字段对应的数据库表的列名，优先级高于类json列名设置，当都没有设置是，默认为字段名
var users3Colsv = &users3Cols{

	Name: "NAME",
	Age:  "age",
	Job:  "job",
	Avg:  "avg",
	Addr: "addr",
}

//指定排序字段
type users3Sort struct {
	Age string
	Avg string
}

var user3Sortv = &users3Sort{}

/* 如果没有指定 cidx 默认是非唯一  */
type users3Cidx struct {
	Name              []string //`cidx:"only"`
	AddrName          []string
	JobAddrNameAvg    []string
	JobAddrNameAvgAge []string
	Test              []string
}

//自定义索引
var users3Cidxv = &users3Cidx{
		Name:     []string{"NAME"},
	//	AddrName: []string{"addr", "Name"},
	Test: []string{"NAME", "addr"},
}

//自定义数据库表名
func (u *users3) TblName() string {
	return "users2"
}

func (u *users3) TblType() string {
	return "INCR"
}

func (u *users3) Cols() *users3Cols {
	return users3Colsv
}

func (u *users3) Sort() *users3Sort {
	return user3Sortv
}

func (u *users3) Cidx() *users3Cidx {
	return users3Cidxv
}

func (u *users3) Print2Map() (uinfo map[string]interface{} ){

	uCol := u.Cols()

	uinfo = make(map[string]interface{})

	uinfo[uCol.Addr] = u.Addr
	uinfo[uCol.Name] = u.Name
	uinfo[uCol.Age] = u.Age
	uinfo[uCol.Avg] = u.Avg
	uinfo[uCol.Job] = u.Job

	return
}

// "age,omitempty"
func ParseTags(s string) (string, []string) {
	if len(s) == 0 {
		return "", nil
	}
	sl := strings.Split(s, ",")
	return sl[0], sl[1:]
}

func SplitTags(s string) []string {
	if len(s) == 0 {
		return nil
	}
	return strings.Split(s, ",")
}

/*****
对信息不全

*****/
var db *gossdb.Client

func main() {
	var uu users3
	//var uu2 []users3

	//设置SSDB数据库，IP与端口，默认127.0.0.1 8888
	if err := ssdb.Start( &conf.Config{
		Host:"127.0.0.1",
		Port:8888,
	}); err != nil {
		fmt.Println(err)
		return
	}
	defer ssdb.Close()

	db, err := ssdb.Client()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	//设置SSDB数据库操作对象的表名，表结构，排序字段，索引字段
	_, err = db.TblInit(&uu, uu.Cols(), uu.Sort(), uu.Cidx())
	if err != nil {
		fmt.Println(err)
	}

	//测试补全命令
	uu = users3{
		Name:"enter lu0",
		Addr:"si chuan cheng du5",
		Age:18,
		Job:"soft development7",
		Avg:18328020553,
	}



	//uu.Id = "0000000000000011"
	//uu.Name = "enter lu11"
	//err = db.Ht_update(uu.TblName(),"cidName1",uu.Id,uu.Print2Map())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//调用提供的save接口，对指定对象的数据进行保存
	//info,err := db.Table2(&uu).Save()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(info,uu.Id)


	//err = db.Ht_delete(uu.TblName(),"cidName1",uu.Id,"addr")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//ret,err := db.Ht_row_clear(uu.TblName(),"cidName1","0000000000000005")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(ret)
	//col,sort,cidx,err := db.Ht_clear(uu.TblName(),"cidName1")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(col,sort,cidx)
	//
	//ret,err := db.Ht_col_count(uu.TblName(),"cidName1","0000000000000010")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(ret)
	//val,err := db.Ht_col_get(uu.TblName(),"cidName1","0000000000000010","NAME")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(val)
	//vals,err := db.Multi_ht_col_get(uu.TblName(),"cidName1","0000000000000010","NAME","age","job")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(vals)
	//retmulti,err := db.Multi_ht_get(uu.TblName(),"cidName1",[]string{"NAME","age","job"},
	//[]string{"0000000000000010","0000000000000011","0000000000000012"})
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retmulti)
	//retmultis,err := db.Multi_ht_get_slice(uu.TblName(),"cidName1",[]string{"NAME","age","job"},
	//[]string{"0000000000000010","0000000000000011","0000000000000012"})
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retmultis)
	//retinfo,err := db.Ht_col_scan(uu.TblName(),"cidName1","0000000000000011","addr","job",5)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfo)

	//删除测试数据
	//ret,err := db.Tclear(uu.TblName())
	//if err != nil {
	//	fmt.Println("db Tbl_drop err :",err)
	//}
	//fmt.Println(ret)
	//
	//reti,err := db.Tbl_drop(uu.TblName())
	//if err != nil {
	//	fmt.Println("db Tbl_drop err :",err)
	//}
	//fmt.Println(reti)
	//
	////Tbl_erase 删除指定表的表结构
	//reti,err = db.Tbl_erase(uu.TblName())
	//if err != nil {
	//	fmt.Println("db Tbl_erase err :",err)
	//}
	//fmt.Println(reti)
	//tblScaninfo,err := db.Tbl_scan("","",10)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(tblScaninfo)
	//for k,v := range tblScaninfo{
	//	fmt.Println(k,"---",v)
	//}
	//tblScaninfo,err = db.Tbl_rscan("","",10)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(tblScaninfo)
	//for k,v := range tblScaninfo{
	//	fmt.Println(k,"---",v)
	//}
	//
	//tbllist,err := db.Tbl_list(6)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(tbllist)
	//colscaninfo,err := db.Tbl_col_scan(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(colscaninfo)
	//for k,v := range colscaninfo{
	//	fmt.Println(k,"---",v)
	//}
	//collistinfo,err := db.Tbl_col_list(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(collistinfo)
	//for k,v := range colscaninfo{
	//	fmt.Println(k,"---",v)
	//}
	//colsortinfo,err := db.Tbl_sort_scan(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(colsortinfo)
	//colsortinfo,err = db.Tbl_sort_list(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(colsortinfo)
	//colsortinfo,err = db.Tbl_cidx_scan(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(colsortinfo)
	//colsortinfo,err = db.Tbl_cidx_list(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(colsortinfo)

	//fmt.Println("-------------Tinsert--------------")
	//增加
	//fmt.Println(uu)
	//	//id,err := db.Tinsert(uu.TblName(),"",uu.Print2Map())
	//	//if err != nil{
	//	//	fmt.Println(err)
	//	//}
	//	//fmt.Println(id)
	//修改
	//fmt.Println("-------------Tupdate--------------")
	//uu.Id = id
	////修改字段内容
	//uu.Name = "enter_lu3"
	//uu.Age = 20
	//fmt.Println(uu)
	//err = db.Tupdate(uu.TblName(),uu.Id,uu.Print2Map())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//删除
	////删除字段内容
	//fmt.Println("-------------Tupdate--------------")
	//uu.Id = id
	//fmt.Println(uu)
	//err = db.Tdelete(uu.TblName(),uu.Id,"addr","avg")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//删除一行数据
	//fmt.Println("-------------Trow_clear--------------")
	//uu.Id = id
	//fmt.Println(uu)
	//err = db.Trow_clear(uu.TblName(),uu.Id)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//清空指定表数据
	//fmt.Println("-------------Tclear--------------")
	//retinfo,err := db.Tclear("users")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfo)
	//获取
	//获取指定表ID的列数
	//fmt.Println("-------------Tcol_count--------------")
	//retinfo,err := db.Tcol_count(uu.TblName(),uu.Id)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfo)
	//获取指定表ID的某列的数值
	//fmt.Println("-------------Tcol_get--------------")
	//uu.Id = "0000000000000014"
	//retinfo,err := db.Tcol_get(uu.TblName(),uu.Id,"name")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfo)
	////获取指定表ID的多列的数值
	//fmt.Println("-------------Multi_tcol_get--------------")
	//uu.Id = "0000000000000008"
	//retinfos,err := db.Multi_tcol_get(uu.TblName(),uu.Id,"NAME","age","avg")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfos)
	////获取指定表ID的多列的数值
	//fmt.Println("-------------Multi_tget--------------")
	//retinfos,err := db.Multi_tget(uu.TblName(),[]string{
	//	"0000000000000006","0000000000000001",
	//	"0000000000000008"},[]string{
	//		"NAME",
	//		"age",
	//		"avg"})
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfos)
	//fmt.Println("-------------Multi_tget--------------")
	//retinfos2,err := db.Multi_tget_slice(uu.TblName(),[]string{
	//	"0000000000000006","0000000000000001",
	//	"0000000000000008"},[]string{
	//	"NAME",
	//	"age",
	//	"avg"})
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfos2)
	////获取指定表ID的范围列的对应的数值
	//uu.Id = "000000000000000c"
	//fmt.Println("-------------Tcol_scan--------------")
	//retminfos,err := db.Tcol_scan(uu.TblName(),uu.Id,"addr","avg",5)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retminfos)
	////获取指定表ID的范围列的对应的数值
	//fmt.Println("-------------Tcol_rscan--------------")
	//retminfos,err = db.Tcol_rscan(uu.TblName(),uu.Id,"age","NAME",5)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retminfos)
	//
	////获取指定表ID的范围列的对应的数值
	//fmt.Println("-------------Trow_scan--------------")
	//retinfos,err := db.Trow_scan("users2","","",5,"name","age","avg")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfos)
	//fmt.Println("-------------Trow_scan--------------")
	//retinfoss,err := db.Trow_scanSlice("users2","","",5,"name","age","avg")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfoss)
	//fmt.Println("-------------Trow_rscan--------------")
	//retinfos,err = db.Trow_rscan("users2","","",5,"NAME","age","avg")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfos)
	//fmt.Println("-------------Trow_rscan--------------")
	//retinfoss,err = db.Trow_rscanSlice("users2","","",5,"NAME","age","avg")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(retinfoss)
	////排序 ——升序
	//fmt.Println("-------------------------------------------")
	//num,retSortInfo,err := db.Tsort_scan("users2","age","","",0,5,"NAME")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//fmt.Println(retSortInfo)
	//fmt.Println("-------------------------------------------")
	//
	//num,retSortSliceInfo,err := db.Tsort_scan_slice(uu.TblName(),"age","","",0,5,"*")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//for _,v := range retSortSliceInfo{
	//	fmt.Println(v.Idx,v.Score,v.ColNum,v.Cols,v.Cold)
	//}
	//fmt.Println("---------------- down sort --------------------")
	////排序 ——降序
	//num,retSortInfo,err = db.Tsort_rscan(uu.TblName(),"age","","",0,5,"NAME")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//fmt.Println(retSortInfo)
	//
	//num,retSortSliceInfo,err = db.Tsort_rscan_slice(uu.TblName(),"age","","",0,5,"*")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//for _,v := range retSortSliceInfo{
	//	fmt.Println(v.Idx,v.Score,v.ColNum,v.Cols,v.Cold)
	//}
	//fmt.Println("---------------- up cidx --------------------")
	////索引——升序
	//num,retCidxInfo,err := db.Tcidx_scan(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"NAME",Start_:"enter lu",End_:"enter lu3"}},10,"name","age","addr")
	////num,retCidxInfo,err := db.Tcidx_scan(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"name",Start_:"",End_:""}},10,"age")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//fmt.Println(retCidxInfo)
	//num,retCidxSliceInfo,err := db.Tcidx_scan_slice(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"NAME",Start_:"enter lu2",End_:"enter_lu3"}},10,"name","age","addr")
	////num,retCidxInfo,err := db.Tcidx_scan(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"name",Start_:"",End_:""}},10,"age")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//for _,v := range retCidxSliceInfo{
	//	fmt.Println(v.Idx,v.Cidx,v.ColNum,v.Cols,v.Cold)
	//}
	//fmt.Println("---------------- down cidx --------------------")
	////索引——降序
	////num,retCidxInfo,err := db.Tcidx_rscan(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"name",Start_:"enter lu3",End_:"enter_lu2"}},10,"name","age","addr")
	//num,retCidxInfo,err = db.Tcidx_rscan(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"NAME",Start_:"",End_:""}},10,"age")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//fmt.Println(retCidxInfo)
	//
	//num,retCidxSliceInfo,err = db.Tcidx_rscan_slice(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"NAME",Start_:"",End_:""}},10,"*")
	////num,retCidxInfo,err := db.Tcidx_scan(uu.TblName(),[]gossdb.Cidx_scan_t{{Col:"name",Start_:"",End_:""}},10,"age")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(num)
	//for _,v := range retCidxSliceInfo{
	//	fmt.Println(v.Idx,v.Cidx,v.ColNum,v.Cols,v.Cold)
	//}
	//add,del,err := db.Tsort_reset(uu.TblName(),"age")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(add,del)
	//
	//add,del,err = db.Tcidx_reset(uu.TblName(),"name","addr")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(add,del)
	//ret,err := db.Trow_count(uu.TblName())
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(ret)


	/*
		db.Table(&uu).Get("id")
		db.Table(&uu).Sort(uu.Sort().Avg,32,10)
		db.Table(&uu).Cidx(uu.Cidx().Name,)
		db.Table(&uu).("id")
		Fields("id,name,age").
			Where("id", 55).
			Select()


			tt := uu.Cols()
			fmt.Println("Name:", tt.Name)
			fmt.Println("Age:", tt.Age)
			fmt.Println("Avg:", tt.Avg)
			fmt.Println("Job:", tt.Job)
			fmt.Println("Addr:", tt.Addr)
	*/
}
func test(uu users3) {
	//测试数据库接口命令
	//Tbl_info 查看指定表的表结构
	//info, err := db.Tbl_info(uu.TblName())
	//if err != nil {
	//	fmt.Println("db Tbl_info err :", err)
	//}
	//fmt.Println(info)

	//Tbl_cidxdel 删除表索引字段
	//ret,err := db.Tbl_cidxdel(uu.TblName(),"addr","name")
	//if err != nil {
	//	fmt.Println("db Tbl_sortdel err :",err)
	//}
	//fmt.Println(ret)

	////Tbl_sortdel 删除表排序字段
	//err = db.Tbl_sortdel(uu.TblName(),"age")
	//if err != nil {
	//	fmt.Println("db Tbl_sortdel err :",err)
	//}

	////Tbl_coldel 删除表字段
	//ret,err := db.Tbl_coldel(uu.TblName(),"name","age")
	//if err != nil {
	//	fmt.Println("db Tbl_coldel err :",err)
	//}
	//fmt.Println(ret)

	//Tbl_drop 删除指定表，表的数据不清除
	//ret,err := db.Tbl_drop(uu.TblName())
	//if err != nil {
	//	fmt.Println("db Tbl_drop err :",err)
	//}
	//fmt.Println(ret)
	//
	////Tbl_erase 删除指定表的表结构
	//ret,err = db.Tbl_erase(uu.TblName())
	//if err != nil {
	//	fmt.Println("db Tbl_erase err :",err)
	//}
	//fmt.Println(ret)
}
/*
type TbUserInfo struct {
	Id       int64       `xorm:"pk autoincr unique BIGINT" json:"id"`
	Phone    string      `xorm:"not null unique VARCHAR(20)" json:"phone"`
	UserName string      `xorm:"VARCHAR(20)" json:"user_name"`
	Gender   int         `xorm:"default 0 INTEGER" json:"gender"`
	Pw       string      `xorm:"VARCHAR(100)" json:"pw"`
	Token    string      `xorm:"TEXT" json:"token"`
	Avatar   string      `xorm:"TEXT" json:"avatar"`
	Extras   interface{} `xorm:"JSON" json:"extras"`
	Created  time.Time   `xorm:"DATETIME created"`
	Updated  time.Time   `xorm:"DATETIME updated"`
	Deleted  time.Time   `xorm:"DATETIME deleted"`
}
*/
