package main

import (
	"fmt"
	"github.com/seefan/gossdb"
	"github.com/seefan/gossdb/conf"
	"github.com/seefan/sorm"
	"github.com/seefan/sorm/example/test/data"
)
var pool *gossdb.Connectors = nil

func main() {

	//1.初始化
	//配置ssdb数据库客户端环境  返回一个句柄

	//2.根据句柄提供的接口 将传入的类/对象同步到ssdb中的一张表中

	//3.根据句柄提供的接口 对数据进行增删改查

	var db *gossdb.Client
	var err error
	pool,err = gossdb.NewPool(&conf.Config{
		Host:             "127.0.0.1",
		Port:             8888,
		SlaveHost:		  "127.0.0.1",//从服务器地址
		SlavePort:		  8889,//从服务器端口
		RetryEnabled:	  true,//是否重连
		MinPoolSize:      10,
		MaxPoolSize:      100,
		MaxWaitSize:      10000,
		AcquireIncrement: 5,
	})
	if err != nil {
		fmt.Println("create pool error", err)
		return
	}
	defer pool.Close()
	db,err = pool.NewClient()
	if err != nil {
		fmt.Println("create new client error", err)
		return
	}
	defer db.Close()

	////3.调用orm对应接口：
	//
	////3.1获取对象对应的ssdb数据库表信息
	//ret1,err := orm.Info()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println("+++++++++++++++++++++++++++++++++++++++")
	//	fmt.Println("+ cols name:",ret1.ColName)
	//	fmt.Println("+ cols type:",ret1.ColType)
	//	fmt.Println("+++++++++++++++++++++++")
	//	fmt.Println("+ sort type:",ret1.SortCols)
	//	fmt.Println("+++++++++++++++++++++++")
	//	fmt.Println("+ cidx name:",ret1.CidxCols)
	//	fmt.Println("+ cidx type:",ret1.CidxType)
	//	fmt.Println("+++++++++++++++++++++++++++++++++++++++")
	//
	//}
	//3.2隐藏表结构
	//ret2,err := orm.Drop()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(ret2)
	//}
	////3.3清除表结构
	//ret3,err := orm.Erase()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(ret3)
	//}
	//3.4.1列出表的列信息 列名升序
	//ret41,err := orm.Scan_colName()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	for k,v := range ret41{
	//		fmt.Print(k,v,"---")
	//	}
	//	fmt.Println()
	//}
	////3.4.2列出表的列信息 列名无序
	//ret42,err := orm.List_colName()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	for k,v := range ret42{
	//		fmt.Print(k,v,"---")
	//	}
	//	fmt.Println()
	//}
	//////3.5.1列出表的排序字段列信息 列名升序
	//ret51,err := orm.Scan_sortName()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	for k,v := range ret51{
	//		fmt.Print(k,v,"---")
	//	}
	//	fmt.Println()
	//}
	//////3.5.2列出表的排序字段列信息 列名无序
	//ret52,err := orm.List_sortName()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	for k,v := range ret52{
	//		fmt.Print(k,v,"---")
	//	}
	//	fmt.Println()
	//}
	//////3.6.1列出表的索引字段列信息 列名升序
	//ret61,err := orm.Scan_cidxName()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	for k,v := range ret61{
	//		fmt.Print(k,v,"---")
	//	}
	//	fmt.Println()
	//}
	//////3.6.2列出表的索引字段列信息 列名无序
	//ret62,err := orm.List_cidxName()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	for k,v := range ret62{
	//		fmt.Print(k,v,"---")
	//	}
	//	fmt.Println()
	//}
	//3.7根据对象数据信息增加一行数据到对应的ssdb数据table中
	//uu := data.Users1{
	//	//Id:"123460",
	//	Name:"wan1gw23u5",
	//	Age:220,
	//	Avg:1122334455,
	//	Sex:false,
	//	Job:"soft development",
	//	Addr:"si chuan cheng du2",
	//}
	//
	//uu.Name = "按时大大"
	//uu.Sex = true
	//id ,err := uu.Add("2222",db)
	//id ,err := orm.Tinsert(&uu)
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(id)
	//}
	uu2 := data.Users1{}
	uu2.Id = "4097"
	err = sorm.CreateTables(db, uu2.TblName())
	orm := sorm.Engine{DB:db}
	orm.SetTable(sorm.Tables[uu2.TblName()])
	ids,err := orm.Htinsert("2222",uu2)
	//_,err = uu2.Get("2222",db,uu2.Cols().Name,uu2.Cols().Age)
	if err != nil {
		fmt.Println(err)
	}else{
		//fmt.Printf(" uu name : %s-%d\n",uu2.Name,uu2.Age)
		fmt.Println(ids)
	}

	//num,ret,err := orm.Rscan_ht_cidx("2222",gossdb.Cidx_scan_t{uu.Cidx().Age,0,nil,"",""},5,uu.Cols().Sex)
	//if err != nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(ret)
	//}
	//ret,err := orm.Scan_trow("","",4,"*")
	//for _,v := range ret{
	//	fmt.Printf("[%s] %v\n",v.Id,v.Colv)
	//}
	/****************************  命令行功能测试  ***********************************/

	/*
	var opt,limit int
	var cols string

	for ;;{
opt:
		fmt.Println("***************************************")
		fmt.Println("*\t 1.add")
		fmt.Println("*\t 2.find")
		fmt.Println("*\t 3.update")
		fmt.Println("*\t 4.delete")
		fmt.Println("*\t 0.quit")
		fmt.Println("***************************************")
		fmt.Println("please input :")
		fmt.Scanf("%d",&opt)
		switch opt {
		case 1:fmt.Println("input add    information:")
			fmt.Println("idx:")
			fmt.Scanf("%s",&uu.Id)
			fmt.Println("Name:")
			fmt.Scanf("%s",&uu.Name)
			fmt.Println("age:")
			fmt.Scanf("%d",&uu.Age)
			fmt.Println("avg:")
			fmt.Scanf("%d",&uu.Avg)
			fmt.Println("Job:")
			fmt.Scanf("%s",&uu.Job)
			fmt.Println("Addr:")
			fmt.Scanf("%s",&uu.Addr)

			ret,err := orm.Tinsert(uu)
			if err != nil{
				fmt.Println(err)
			}else {
				fmt.Println(ret)
			}
		case 2:fmt.Println("input find   information:")
			fmt.Println("find ids:")
			fmt.Scanf("%s",&uu.Id)
			fmt.Println("find limit:")
			fmt.Scanf("%d",&limit)
			fmt.Println("find cols:")
			fmt.Scanf("%s",&cols)
			ret,err := orm.Scan_trow("","",limit,cols)
			if err != nil{
				fmt.Println(err)
			}else {
				if uu.Id == ""{
					fmt.Println(ret)
				}else{
					for _,v :=range ret{
						if v.Id == uu.Id{
							fmt.Println(v)
							break
						}
					}
				}


			}
		case 3:fmt.Println("input update information:")
		case 4:fmt.Println("input delete information:")
		case 0:
			return
		default:
			fmt.Println("opt arg set error please reset...")
			goto opt
		}

	}
	*/
	//uu3 := users1{Name:"lulang3"}
	//
	//err = orm.Tupdate(uu2.Id,&uu3)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//
	//err = orm.Tupdate(uu.Id,map[string]interface{}{
	//	"NAME":"lulang4",
	//	"Age":22,
	//})
	//if err != nil{
	//	fmt.Println(err)
	//}

	//err = orm.Tdelete(uu2.Id,"Avg","Addr")
	//if err != nil{
	//	fmt.Println(err)
	//}

	//err = orm.Clear_trow(uu2.Id)
	//if err != nil{
	//	fmt.Println(err)
	//}

	//列数
	//col,err := orm.Count_col(uu.Id)
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(col)
	//}
	////行数
	//trow,err := orm.Count_trow()
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(trow)
	//}
	////指定id,col,获取对应数值
	//val ,err := orm.Get_col(uu.Id,"Name")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println("name:",val)
	//}
	//
	//vals ,err := orm.Multi_Get_cols(uu.Id,"Name","age")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println("Name:",vals[0],"age:",vals[1])
	//}
	//
	//valsm,err := orm.Multi_Gets_cols([]string{uu.Id,uu2.Id},[]string{"Name","age"})
	//
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valsm)
	//}
	//
	//valss,err := orm.Multi_Gets_cols_slice([]string{uu.Id,uu2.Id},[]string{"Name","age"})
	//
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valss)
	//}
	//
	//valmss,err := orm.Scan_col(uu.Id,"","",10)
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valmss)
	//}
	//
	//valrs,err := orm.Scan_trow("","",10,"Name")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valrs)
	//}
	//valrs,err = orm.Rscan_trow("","",10,"*")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valrs)
	//}
	//
	//valrss,err := orm.Scan_trow_slice("","",10,"*")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valrss)
	//}
	//valrss,err = orm.Rscan_trow_slice("","",10,"*")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(valrss)
	//}
	//num,valsss,err := orm.Scan_sort("age","","",0,10,"Name")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valsss)
	//}
	//num,valsss,err = orm.Rscan_sort("age","","",0,10,"Name")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valsss)
	//}
	//num,valssss,err := orm.Scan_sort_slice("age","","",0,10,"Name")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valssss)
	//}
	//num,valssss,err = orm.Rscan_sort_slice("age","","",0,10,"Name")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valssss)
	//}
	//cidxs := []gossdb.Cidx_scan_t{{"Name","",""}}
	//num,valsc,err := orm.Scan_cidx(cidxs,10,"Name","Age","Addr")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valsc)
	//}
	//num,valsc,err = orm.Rscan_cidx(cidxs,10,"Name","Age","Addr")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valsc)
	//}
	//num,valscs,err := orm.Scan_cidx_slice(cidxs,10,"Name","Age","Addr")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valscs)
	//}
	//num,valscs,err = orm.Rscan_cidx_slice(cidxs,10,"Name","Age","Addr")
	//if err != nil{
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(num)
	//	fmt.Println(valscs)
	//}

	//ret1,ret2,err := orm.Reset_sort("Age")
	//if err != nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(ret1,ret2)
	//}
	//cidxinfo,err := orm.Scan_cidxName()
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(cidxinfo)
	//ret1,ret2,err = orm.Reset_cidx(cidxinfo[0].Cidx...)
	//if err != nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(ret1,ret2)
	//}

	/****************************  htable 集合表  ***********************************/
	//rets,err := orm.Htinsert("test1",uu,uu2)
	//if err != nil {
	//	fmt.Println(rets)
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(rets)
	//}


}
