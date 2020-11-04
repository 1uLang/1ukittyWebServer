package sorm

import (
	"fmt"
	"github.com/seefan/gossdb"
)

type Table struct {
	tblName		string				//表名
	tblType		string				//表类型
	keyName		string				//表主键名
	colsName	map[string]string	//类成员对应列名
	colsType	map[string]string	//表列名对应类型

	sortCols	[]string 			//排序字段名
	sortState   map[string]bool		//方便对排序字段遍历
	cidxCols    map[string][]string //索引名- 对应索引列名[可多个,有序]
	cidxType	map[string]bool		//索引名- 是否唯一
}

var Tables map[string]*Table

func (tbl *Table) Set_tblName(name string)error {
	if tbl == nil{
		return fmt.Errorf("tbl is nil")
	}
	tbl.tblName = name
	return nil
}
func (tbl *Table)check_cols(col ...string)([]string,error){

	if len(col) >0 && col[0] == "*"{
		return col,nil
	}
	cols := make([]string,0)
	for _,v := range col{

		//1.先检测table列名中是否存在改列名
		if tbl.colsType[v] != ""{
			cols = append(cols, v)
		//2.再检测同步的对象中是否存在该字段
		} else if name := tbl.colsName[v]; name != ""{
			cols = append(cols, name)
		} else {
			return nil,fmt.Errorf("%s table not \"%s\" col",tbl.tblName,v)
		}
	}
	return cols,nil
}

func (tbl *Table) check_sort(s string) (string,error) {

	col ,err := tbl.check_cols(s)
	if err != nil{
		return "",err
	}
	if tbl.sortState[col[0]]{
		return col[0],nil
	}

	return "", fmt.Errorf("%s not \"%s\" sort col",tbl.tblName,s)
}

func (tbl *Table) check_cidx(cidx gossdb.Cidx_scan_t) (gossdb.Cidx_scan_t, error) {


	newcidx := cidx

	cols,err := tbl.check_cols(cidx.Cols...)
	if err != nil{
		return cidx,err
	}

	for _,k := range tbl.cidxCols{
		if StringSliceEqual(k,cols){
			return newcidx,nil
		}
	}
	return cidx,fmt.Errorf("%s not \"%v\" cidx cols",tbl.tblName,cidx.Cols)
}

func (tbl *Table) check_cidx2(col ...string) ([]string, error) {

	cols,err := tbl.check_cols(col...)
	if err != nil{
		return nil,err
	}

	for _,k := range tbl.cidxCols{
		if StringSliceEqual(k,cols){
			return cols,nil
		}
	}
	return nil,fmt.Errorf("%s not \"%v\" cidx cols",tbl.tblName,col)
}
