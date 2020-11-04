package sorm

import (
	"fmt"
	"github.com/seefan/gossdb"
	"github.com/seefan/gossdb/conf"
	"github.com/seefan/gossdb/ssdb"
	"reflect"
	"strings"
)

type Engine struct {
	DB  *gossdb.Client
	table *Table
}

func NewEngine(host string,port int)  (*Engine,error){

	if err := ssdb.Start( &conf.Config{
		Host:host,
		Port:port,
	}); err != nil {
		return nil,err
	}

	db, err := ssdb.Client()
	if err != nil {
		return nil,err
	}

	engine := Engine{
		DB:db,
		table:&Table{},
	}

	return &engine,nil
}

func NewEngineAddPool(db *gossdb.Client)  (*Engine,error){

	if db == nil{
		return nil,fmt.Errorf("set db is nil")
	}

	engine := Engine{
		DB:db,
		table:&Table{},
	}

	return &engine,nil
}
func (engine *Engine) Close() {
	if engine.DB != nil{
		_ =engine.DB.Close()
	}
}

func Sync(obj interface{}) error {
	// 获取表的信息
	var (
		tblName    string
		tblType    = "DATA"
		tblFieds   = map[string]string{} /*存放 struct tag */
		colss      = map[string]string{} /*存放 colname coltype*/
		colvs      = map[string]string{} /*存放 struct colname*/
		sortss     = make([]string,0)
		sortst     = map[string]bool{}
		cidxss     = map[string][]string{}
		cidxssonly = map[string]bool{}
	)
	tbl := &Table{}

	rvp := reflect.ValueOf(obj)
	if rvp.Kind() != reflect.Ptr {
		return fmt.Errorf("tbl not Ptr")
	}

	rv := rvp.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("tbl not struct")
	}

	tblName = rv.Type().Name()

	//同步对象，获取映射的table名与类型
	if tn := rvp.MethodByName("TblName"); tn.IsValid() {
		tblName = tn.Call(nil)[0].String()
	}

	if tn := rvp.MethodByName("TblType"); tn.IsValid() {
		tblType = tn.Call(nil)[0].String()
	}

	tbl.tblName = tblName
	tbl.tblType = tblType

	//同步对象，获取对应各成员字段映射的table 列名与类型,主键行id名
	rvt := rv.Type()
	for i := 0; i < rvt.NumField(); i++ {
		field := rvt.Field(i)

		// unexported field
		if field.PkgPath != "" {
			continue
		}
		tbltags := splitTags(field.Tag.Get("json"))

		if tbltags == nil || tbltags[0] == "-" || tbltags[0] == "omitempty" {

			tblFieds[field.Name] = field.Name
		}else {
			tblFieds[field.Name] = tbltags[0]
		}

		if strings.Contains(string(field.Tag),"tbl:\"KEY\""){
			tbl.keyName = field.Name
		}
	}

	//获取 对象设置table列名以及对应类型信息
	if tn := rvp.MethodByName("Cols"); tn.IsValid() {
		rvp := tn.Call(nil)[0]
		if rvp.Kind() != reflect.Ptr {
			return fmt.Errorf("cols not Ptr")
		}
		rv = rvp.Elem()
		if rv.Kind() != reflect.Struct {
			return fmt.Errorf("cols not struct")
		}

		rvt = rv.Type()
		for i := 0; i < rvt.NumField(); i++ {
			field := rvt.Field(i)
			// unexported field
			if field.PkgPath != "" {
				continue
			}
			colType := "DATA"
			fieldValue := rv.FieldByName(field.Name)
			colName := fieldValue.Interface().(string)
			if colName == "" {
				colName = field.Name
				if tt, rr := tblFieds[colName]; rr {
					if tt != "" {
						colName = tt
					}
				}
				fieldValue.SetString(colName)
			}

			tag := field.Tag.Get("tbl")
			if tag != "" {
				colType = tag
			}

			colss[colName] = colType
			colvs[field.Name] = colName
		}
	}
	tbl.colsName = colvs
	tbl.colsType = colss

	//获取 对象设置table排序字段列名
	if tn := rvp.MethodByName("Sort"); tn.IsValid() {
		rvp := tn.Call(nil)[0]
		if rvp.Kind() != reflect.Ptr {
			return fmt.Errorf("sorts not Ptr")
		}
		rv = rvp.Elem()
		if rv.Kind() != reflect.Struct {
			return fmt.Errorf("sorts not struct")
		}
		rvt := rv.Type()
		for i := 0; i < rvt.NumField(); i++ {
			field := rvt.Field(i)
			if field.PkgPath != "" {
				continue
			}
			vv, ok := colvs[field.Name]
			if !ok {
				return fmt.Errorf("sorts:'%s' not cols", field.Name)
			}
			/*类型判断*/
			if ctype, ok := colss[vv]; !ok || (ctype != "NUM" && ctype != "!NUM") {
				return fmt.Errorf("sorts:'%s' colType Not NUM", field.Name)
			}
			fieldValue := rv.FieldByName(field.Name)
			fieldValue.SetString(vv)
			sortss = append(sortss, vv)
			sortst[vv] = true
		}
	}
	tbl.sortCols = sortss
	tbl.sortState = sortst


	//获取 对象设置table索引字段列名 - 类型
	if tn := rvp.MethodByName("Cidx"); tn.IsValid() {
		rvp := tn.Call(nil)[0]
		if rvp.Kind() != reflect.Ptr {
			return fmt.Errorf("cidxs not Ptr")
		}
		rv = rvp.Elem()
		if rv.Kind() != reflect.Struct {
			return fmt.Errorf("cidxs not struct")
		}
		rvt := rv.Type()
		for i := 0; i < rvt.NumField(); i++ {
			field := rvt.Field(i)
			if field.PkgPath != "" {
				continue
			}
			fieldValue := rv.FieldByName(field.Name)
			vals := fieldValue.Interface().([]string)
			if len(vals) == 0 {
				vvv := make([]string,0)
				tt := field.Name
				for tt != "" {
					found := false
					for k, v := range colvs {
						if pos := strings.Index(tt, k); pos == 0 {
							vvv = append(vvv, v)
							tt = tt[len(k):]
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("cidxs:'%s' in '%s' name err", field.Name, tt)
					}
				}
				fieldValue.Set(reflect.ValueOf(vvv))
			} else {
				for _, v := range vals {
					if _, ok := colss[v]; !ok {
						return fmt.Errorf("cidxs:'%s' in '%s' not cols", field.Name, v)
					}
				}
			}
			cidxssonly[field.Name] = false
			tags := field.Tag.Get("cidx")
			if tags == "only" {
				cidxssonly[field.Name] = true
			}
			vals = fieldValue.Interface().([]string)

			cidxss[field.Name] = vals
		}
	}
	tbl.cidxCols = cidxss
	tbl.cidxType = cidxssonly

	if Tables == nil{
		Tables = make(map[string]*Table)
	}

	Tables[tblName] = tbl
	return nil
}
func (engine *Engine) SetTable(tbl *Table) {
	engine.table = tbl
}

func CreateTables(db *gossdb.Client,tblName ...string) error {
	c := db

	if c == nil {
		return fmt.Errorf("engine DB is nil")
	}
	engine := Engine{DB:db}
	msg := ""
	fmt.Println("create tables tbl name : ",tblName)
	for _,v := range tblName{
		tbl := Tables[v]
		engine.table = tbl
		if tbl == nil {
			msg += "tbl %s is not exists\n"
			continue
		}
		//创建表
		state,err := c.Tbl_create(tbl.tblName, tbl.tblType)
		if err != nil{
			return err
		}

		if state == 0{	//已存在
			return nil
		}

		//添加字段
		_, err = c.Tbl_coladd(tbl.tblName, tbl.colsType)
		if err != nil {
			return err
		}

		//添加排序
		for _, k := range tbl.sortCols{
			_, err := c.Tbl_sortadd(tbl.tblName, k)
			if err != nil {
				return err
			}
		}

		//添加索引

		for k, cidx := range tbl.cidxCols {
			var cc []string
			cc = append(cc, cidx...)
			if tbl.cidxType[k] {
				cc = append(cc, "1")
			} else {
				cc = append(cc, "0")
			}
			err := c.Tbl_cidxadd(tbl.tblName, cc...)
			if err != nil {
				return err
			}
		}
		//reset cidx sort col
		for _,v := range tbl.cidxCols{
			_,_,err = engine.Reset_cidx(v...)
			if err != nil {
				return fmt.Errorf("tbl reset cidx err:%s",err.Error())
			}
		}
		for _,v := range tbl.sortCols{
			_,_,err = engine.Reset_sort(v)
			if err != nil {
				return fmt.Errorf("tbl reset sort err:%s",err.Error())
			}
		}

	}
	if msg != ""{
		fmt.Errorf(msg)
	}
	return nil
}

func check_table(tbl *Table)  error{

	if tbl == nil{
		return fmt.Errorf("engine table is nil,you need to Sync obj")
	}
	return nil
}
func (engine *Engine) Info() (*gossdb.Tblinfo, error) {

	if err := check_table(engine.table); err != nil{
		return nil,err
	}

	return engine.DB.Tbl_info(engine.table.tblName)
}

func (engine *Engine) Drop() (int64, error) {

	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Tbl_drop(engine.table.tblName)
}

func (engine *Engine) Erase() (int64, error) {

	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Tbl_erase(engine.table.tblName)
}

func (engine *Engine) Scan_tbl(keyStart, keyEnd string, limit int64) (map[string]string, error) {

	return engine.DB.Tbl_scan(keyStart,keyStart,limit)
}

func (engine *Engine) Rscan_tbl(keyStart, keyEnd string, limit int64) (map[string]string, error) {

	return engine.DB.Tbl_rscan(keyStart,keyStart,limit)
}

func (engine *Engine) List_tbl(limit int) (map[string]string, error) {

	return engine.DB.Tbl_list(limit)
}

func (engine *Engine) Scan_colName() (map[string]string, error) {

	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tbl_col_scan(engine.table.tblName)
}
func (engine *Engine) List_colName() (map[string]string, error) {

	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tbl_col_list(engine.table.tblName)
}
func (engine *Engine) Scan_sortName() ([]string, error) {

	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tbl_sort_scan(engine.table.tblName)
}
func (engine *Engine) List_sortName() ([]string, error){

	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tbl_sort_list(engine.table.tblName)
}

func (engine *Engine) Scan_cidxName() ([]gossdb.CidxInfo, error) {

	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tbl_cidx_scan(engine.table.tblName)
}
func (engine *Engine) List_cidxName() ([]string, error){

	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tbl_cidx_list(engine.table.tblName)
}
/****************************  table 集合表  ***********************************/
func (engine *Engine) Tinsert(beans ...interface{}) (ids []string,err error){

	c := engine.DB
	tbl := engine.table
	ids = make([]string,0)

	if err := check_table(tbl); err != nil{
		return nil,err
	}
	//当添加失败，将原添加的数据删除清空
	defer func() {
		if err != nil && ids != nil && len(ids) != 0{
			fmt.Println(ids)
			for _,v := range ids{
				_,_ = engine.Clear_trow(v)
			}
			ids = nil
		}
	}()

	for _, bean := range beans {
		switch bean.(type) {
		case map[string]interface{}:
			id, err := c.Tinsert(tbl.tblName,"",bean.(map[string]interface{}))
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		case []map[string]interface{}:
			s := bean.([]map[string]interface{})
			for i := 0; i < len(s); i++ {
				id, err := c.Tinsert(tbl.tblName,"",bean.(map[string]interface{}))
				if err != nil {
					return ids, err
				}
				ids = append(ids, id)
			}
		case map[string]string:
			id, err := c.TinsertMapString(tbl.tblName,"",bean.(map[string]string))
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		case []map[string]string:
			s := bean.([]map[string]interface{})
			for i := 0; i < len(s); i++ {
				id, err := c.TinsertMapString(tbl.tblName,"",bean.(map[string]string))
				if err != nil {
					return ids, err
				}
				ids = append(ids, id)
			}
		default:
			rvp := reflect.ValueOf(bean)
			var rv reflect.Value
			colval := make(map[string]interface{},0)
			if rvp.Kind() == reflect.Ptr {
				rv = rvp.Elem()
			}else {
				rv = rvp
			}

			if rv.Kind() != reflect.Struct {
				return ids, fmt.Errorf("tbl not struct")
			}
			rvt := rv.Type()
			id := ""
			for i := 0; i < rvt.NumField(); i++ {
				field := rvt.Field(i)
				// unexported field
				if field.PkgPath != "" {
					continue
				}
				fieldValue := rv.FieldByName(field.Name)

				if tbl.keyName == field.Name {
					id = fieldValue.String()

				}
				if tbl.colsName[field.Name] == ""{
					continue
				}
				if fieldValue.String() != "" {
					colval[tbl.colsName[field.Name]] = fieldValue.Interface()
				}
			}
			id,err := c.Tinsert(tbl.tblName,id,colval)
			if err != nil{
				return ids,err
			}
			ids = append(ids, id)
		}
	}
	return ids,nil
}

func isZeroData(data reflect.Value) bool{

	switch data.Kind() {
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
			return data.Int() == 0
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			return data.Uint() == 0
	case reflect.String:
			return data.String() == ""
	default:
		return data.IsValid()
	}
}
func (engine *Engine) Tupdate(id string,bean interface{})error{

	if err := check_table(engine.table); err != nil{
		return err
	}
	if id == "" {
		return fmt.Errorf("condition key id is can not null")
	}
	tbl := engine.table
	r := reflect.ValueOf(bean)
	rt := r.Kind()
	cv := make(map[string]interface{})

	if rt == reflect.Struct|| (rt == reflect.Ptr && r.Elem().Kind() == reflect.Struct){
		if rt == reflect.Ptr{
			r = r.Elem()
		}
		for k,v := range tbl.colsName{
			val := r.FieldByName(k)
			if val.IsValid() && val.String() != ""{
				cv[v] = val.Interface()
			}
		}
	}else if rt == reflect.Map{
			ks := r.MapKeys()
			for _,v := range ks{
				if v.Kind() != reflect.String{
					return fmt.Errorf("args bean key type err not string")
				}
				//检测设置的更新数据中,列名是否存在
				cols,err := tbl.check_cols(v.String())
				if err != nil{
					return err
				}
				cv[cols[0]] = r.MapIndex(v).Interface()
			}
	}else {
		return fmt.Errorf("args bean is not Stru,Ptr or Map")
	}

	return engine.DB.Tupdate(tbl.tblName,id,cv)
}

func (engine *Engine) Tdelete(id string,col ...string)error{
	if err := check_table(engine.table); err != nil{
		return err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(col...)
	if err != nil{
		return  err
	}
	return engine.DB.Tdelete(engine.table.tblName,id,cols...)
}
func (engine *Engine)Clear_trow(id string)(int,error)  {
	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Trow_clear(engine.table.tblName,id)
}
func (engine *Engine)Clear_table()([]string,error ) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tclear(engine.table.tblName)
}
func (engine *Engine)Count_col(id string)(int64,error ) {
	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Tcol_count(engine.table.tblName,id)
}
func (engine *Engine)Count_trow()(int64,error ) {
	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Trow_count(engine.table.tblName)
}
func (engine *Engine)Get_col(id,col string)(gossdb.Value,error ) {
	if err := check_table(engine.table); err != nil{
		return "",err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(col)
	if err != nil{
		return "", err
	}
	return engine.DB.Tcol_get(engine.table.tblName,id,cols[0])
}
func (engine *Engine)Multi_Get_cols(id string,col ...string)([]gossdb.Value,error ) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(col...)
	if err != nil{
		return nil, err
	}

	return engine.DB.Multi_tcol_get(engine.table.tblName,id,cols...)
}
func (engine *Engine)Multi_Gets_cols(ids ,cols []string)([]gossdb.MultiGet,error ) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(cols...)
	if err != nil{
		return nil, err
	}
	return engine.DB.Multi_tget(engine.table.tblName,ids,cols)
}
func (engine *Engine)Multi_Gets_cols_slice(ids ,cols []string)([]gossdb.MultiGetSlice,error ) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(cols...)
	if err != nil{
		return nil, err
	}
	return engine.DB.Multi_tget_slice(engine.table.tblName,ids,cols)
}
func (engine *Engine) Scan_col(id, col_start, col_end string, limit int)(map[string]gossdb.Value, error)   {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tcol_scan(engine.table.tblName,id,col_start,col_end,limit)
}
func (engine *Engine) Rscan_col(id, col_start, col_end string, limit int)(map[string]gossdb.Value, error)   {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Tcol_rscan(engine.table.tblName,id,col_start,col_end,limit)
}

func (engine *Engine) Scan_trow(id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScan, error) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(cols...)
	if err != nil{
		return nil, err
	}
	return engine.DB.Trow_scan(engine.table.tblName,id_start,id_end,limit,cols...)
}
func (engine *Engine) Rscan_trow(id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScan, error) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(cols...)
	if err != nil{
		return nil, err
	}
	return engine.DB.Trow_rscan(engine.table.tblName,id_start,id_end,limit,cols...)
}
func (engine *Engine) Scan_trow_slice(id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScanSlice, error) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(cols...)
	if err != nil{
		return nil, err
	}
	return engine.DB.Trow_scan_slice(engine.table.tblName,id_start,id_end,limit,cols...)
}
func (engine *Engine) Rscan_trow_slice(id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScanSlice, error) {
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err := engine.table.check_cols(cols...)
	if err != nil{
		return nil, err
	}
	return engine.DB.Trow_rscan_slice(engine.table.tblName,id_start,id_end,limit,cols...)
}

func (engine *Engine) Scan_sort(sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScan, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table是否拥有指定的sort排序列名
	sortName,err = engine.table.check_sort(sortName)
	if err != nil{
		return 0, nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tsort_scan(engine.table.tblName,sortName,start_,end_,offset,limit,cols...)
}
func (engine *Engine) Rscan_sort(sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScan, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table是否拥有指定的sort排序列名
	sortName,err = engine.table.check_sort(sortName)
	if err != nil{
		return 0, nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tsort_rscan(engine.table.tblName,sortName,start_,end_,offset,limit,cols...)
}
func (engine *Engine) Scan_sort_slice(sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScanSlice, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table是否拥有指定的sort排序列名
	sortName,err = engine.table.check_sort(sortName)
	if err != nil{
		return 0, nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tsort_scan_slice(engine.table.tblName,sortName,start_,end_,offset,limit,cols...)
}
func (engine *Engine) Rscan_sort_slice(sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScanSlice, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table是否拥有指定的sort排序列名
	sortName,err = engine.table.check_sort(sortName)
	if err != nil{
		return 0, nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tsort_rscan_slice(engine.table.tblName,sortName,start_,end_,offset,limit,cols...)
}

func (engine *Engine)Scan_cidx(cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScan, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table中是否存在指定的索引列
	cidx ,err = engine.table.check_cidx(cidx)
	if err != nil{
		return 0,nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tcidx_scan(engine.table.tblName,cidx,offset,limit,cols...)
}

func (engine *Engine)Rscan_cidx(cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScan, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table中是否存在指定的索引列
	cidx ,err = engine.table.check_cidx(cidx)
	if err != nil{
		return 0,nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tcidx_rscan(engine.table.tblName,cidx,offset,limit,cols...)
}

func (engine *Engine)Scan_cidx_slice(cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScanSlice, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table中是否存在指定的索引列
	cidx ,err = engine.table.check_cidx(cidx)
	if err != nil{
		return 0,nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tcidx_scan_slice(engine.table.tblName,cidx,offset,limit,cols...)
}

func (engine *Engine)Rscan_cidx_slice(cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScanSlice, err error) {
	if err := check_table(engine.table); err != nil{
		return 0,nil,err
	}
	//检测table中是否存在指定的索引列
	cidx ,err = engine.table.check_cidx(cidx)
	if err != nil{
		return 0,nil, err
	}
	//检测table，或同步的对象是否拥有指定的列名或对象字段名，返回对用的映射的table列名
	cols,err = engine.table.check_cols(cols...)
	if err != nil{
		return 0,nil, err
	}
	return engine.DB.Tcidx_rscan_slice(engine.table.tblName,cidx,offset,limit,cols...)
}

func (engine *Engine)Reset_sort(col string)(int64,int64,error){
	if err := check_table(engine.table); err != nil{
		return 0,0,err
	}
	col,err := engine.table.check_sort(col)
	if err != nil{
		return 0,0,err
	}
	return engine.DB.Tsort_reset(engine.table.tblName,col)
}

func (engine *Engine)Reset_cidx(cols ...string)(int64,int64,error){
	if err := check_table(engine.table); err != nil{
		return 0,0,err
	}
	cols,err := engine.table.check_cidx2(cols...)
	if err != nil{
		return 0,0,err
	}
	return engine.DB.Tcidx_reset(engine.table.tblName,cols...)
}
/****************************  htable 集合表  ***********************************/
func (engine * Engine)Htinsert(hashName string,beans ...interface{}) (ids []string,err error) {

	c := engine.DB
	tbl := engine.table
	ids = make([]string,0)

	if err := check_table(tbl); err != nil{
		return nil,err
	}

	//当添加失败，将原添加的数据删除清空
	defer func() {
		if err != nil && ids != nil && len(ids) != 0{
			for _,v := range ids{
				_,_ = engine.Clear_trow(v)
			}
			ids = nil
		}

	}()

	for _, bean := range beans {
		switch bean.(type) {
		case map[string]interface{}:
			id, err := c.Ht_insert(tbl.tblName,hashName,"",bean.(map[string]interface{}))
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		case []map[string]interface{}:
			s := bean.([]map[string]interface{})
			for i := 0; i < len(s); i++ {
				id, err := c.Ht_insert(tbl.tblName,hashName,"",bean.(map[string]interface{}))
				if err != nil {
					return ids, err
				}
				ids = append(ids, id)
			}
		case map[string]string:
			id, err := c.Ht_insertMapString(tbl.tblName,hashName,"",bean.(map[string]string))
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		case []map[string]string:
			s := bean.([]map[string]interface{})
			for i := 0; i < len(s); i++ {
				id, err := c.Ht_insertMapString(tbl.tblName,hashName,"",bean.(map[string]string))
				if err != nil {
					return ids, err
				}
				ids = append(ids, id)
			}
		default:
			rvp := reflect.ValueOf(bean)
			var rv reflect.Value
			colval := make(map[string]interface{},0)
			if rvp.Kind() == reflect.Ptr {
				rv = rvp.Elem()
			}else {
				rv = rvp
			}

			if rv.Kind() != reflect.Struct {
				return ids, fmt.Errorf("tbl not struct")
			}
			rvt := rv.Type()
			id := ""
			for i := 0; i < rvt.NumField(); i++ {
				field := rvt.Field(i)
				// unexported field
				if field.PkgPath != "" {
					continue
				}
				fieldValue := rv.FieldByName(field.Name)

				if tbl.keyName == field.Name {
					id = fieldValue.String()

				}
				if tbl.colsName[field.Name] == ""{
					continue
				}
				if fieldValue.String() != ""{
					colval[tbl.colsName[field.Name]] = fieldValue.Interface()
				}
			}
			id,err := c.Ht_insert(tbl.tblName,hashName,id,colval)
			if err != nil{
				return ids,err
			}
			ids = append(ids, id)
		}
	}
	return ids,nil
}

func (engine *Engine) Htupdate(hashName, id string, bean interface{}) error {

	if err := check_table(engine.table); err != nil{
		return err
	}
	if id == "" {
		return fmt.Errorf("condition key id is can not null")
	}
	tbl := engine.table
	r := reflect.ValueOf(bean)
	rt := r.Kind()
	cv := make(map[string]interface{})

	if rt == reflect.Struct || (rt == reflect.Ptr && r.Elem().Kind() == reflect.Struct){
		if rt == reflect.Ptr{
			r = r.Elem()
		}
		for k,v := range tbl.colsName{
			val := r.FieldByName(k)
			if val.IsValid() && !isZeroData(val){
				cv[v] = val.Interface()
			}
		}
	}else if rt == reflect.Map{
		ks := r.MapKeys()
		for _,v := range ks{
			if v.Kind() != reflect.String{
				return fmt.Errorf("args bean key type err not string")
			}
			//检测设置的更新数据中,列名是否存在
			cols,err := tbl.check_cols(v.String())
			if err != nil{
				return err
			}
			cv[cols[0]] = r.MapIndex(v).Interface()
		}
	}else {
		return fmt.Errorf("args bean is not Stru,Ptr or Map")
	}
	return engine.DB.Ht_update(tbl.tblName,hashName,id,cv)
}
func (engine *Engine) Htdelete(hashName, id string, col ...string) error {
	if err := check_table(engine.table); err != nil{
		return err
	}
	return engine.DB.Ht_delete(engine.table.tblName,hashName,id,col...)
}
func (engine *Engine) Clear_ht_row(hashName, id string)(int64,error){
	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Ht_row_clear(engine.table.tblName,hashName,id)
}
func (engine *Engine) Clear_htable(hashName string)(int64,int64,int64,error){
	if err := check_table(engine.table); err != nil{
		return 0,0,0,err
	}
	return engine.DB.Ht_clear(engine.table.tblName,hashName)
}
func (engine *Engine) Count_ht_col(hashName,id string)(int64,error){
	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Ht_col_count(engine.table.tblName,hashName,id)
}
func (engine *Engine) Count_ht_row(hashName string)(int64,error){
	if err := check_table(engine.table); err != nil{
		return 0,err
	}
	return engine.DB.Ht_row_count(engine.table.tblName,hashName)
}
func (engine *Engine) Get_ht_col(hashName,id,col string)(gossdb.Value,error){
	if err := check_table(engine.table); err != nil{
		return "",err
	}
	return engine.DB.Ht_col_get(engine.table.tblName,hashName,id,col)
}
func (engine *Engine) Multi_get_ht_cols(hashName,id string,cols ...string)([]gossdb.Value,error){
	if err := check_table(engine.table); err != nil{
		return nil,err
	}
	return engine.DB.Multi_ht_col_get(engine.table.tblName,hashName,id,cols...)
}
func (engine *Engine) Multi_gets_ht_cols_slice(hashName string,ids []string,cols []string)([]gossdb.MultiGetSlice,error) {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Multi_ht_get_slice(engine.table.tblName, hashName, ids, cols)
}
func (engine *Engine)Scan_ht_col(hashName,id,col_start,col_end string,limit int)(map[string]gossdb.Value,error)  {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Ht_col_scan(engine.table.tblName, hashName, id, col_start,col_end,limit)
}

func (engine *Engine)Rscan_ht_col(hashName,id,col_start,col_end string,limit int)(map[string]gossdb.Value,error)  {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Ht_col_rscan(engine.table.tblName, hashName, id, col_start,col_end,limit)
}

func (engine *Engine)Scan_ht_row(hashName, id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScan, error) {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Ht_row_scan(engine.table.tblName, hashName, id_start, id_end,limit,cols...)
}
func (engine *Engine)Rscan_ht_row(hashName, id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScan, error) {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Ht_row_rscan(engine.table.tblName, hashName, id_start, id_end,limit,cols...)
}
func (engine *Engine)Scan_ht_row_slice(hashName, id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScanSlice, error) {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Ht_row_scan_slice(engine.table.tblName, hashName, id_start, id_end,limit,cols...)
}
func (engine *Engine)Rscan_ht_row_slice(hashName, id_start, id_end string, limit int ,cols ...string) ([]gossdb.TrowScanSlice, error) {
	if err := check_table(engine.table); err != nil {
		return nil, err
	}
	return engine.DB.Ht_row_rscan_slice(engine.table.tblName, hashName, id_start, id_end,limit,cols...)
}

func (engine *Engine)Scan_ht_sort(hashName, sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScan, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_sort_scan(engine.table.tblName, hashName, sortName,start_, end_,offset,limit,cols...)
}
func (engine *Engine)Rscan_ht_sort(hashName, sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScan, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_sort_rscan(engine.table.tblName, hashName, sortName,start_, end_,offset,limit,cols...)
}
func (engine *Engine)Scan_ht_sort_slice(hashName, sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScanSlice, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_sort_scan_slice(engine.table.tblName, hashName, sortName,start_, end_,offset,limit,cols...)
}
func (engine *Engine)Rscan_ht_sort_slice(hashName, sortName string, start_, end_ interface{},
	offset, limit int64, cols ...string) (num int64, data []gossdb.SortScanSlice, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_sort_rscan_slice(engine.table.tblName, hashName, sortName,start_, end_,offset,limit,cols...)
}

func (engine *Engine)Scan_ht_cidx(hashName string, cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScan, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_cidx_scan(engine.table.tblName, hashName, cidx,offset,limit,cols...)
}
func (engine *Engine)Rscan_ht_cidx(hashName string, cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScan, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_cidx_rscan(engine.table.tblName, hashName, cidx,offset,limit,cols...)
}
func (engine *Engine)Scan_ht_cidx_slice(hashName string, cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScanSlice, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_cidx_scan_slice(engine.table.tblName, hashName, cidx,offset,limit,cols...)
}
func (engine *Engine)Rscan_ht_cidx_slice(hashName string, cidx gossdb.Cidx_scan_t,offset, limit int64,cols ...string) (num int64, data []gossdb.CidxScanSlice, err error) {
	if err := check_table(engine.table); err != nil {
		return 0,nil, err
	}
	return engine.DB.Ht_cidx_rscan_slice(engine.table.tblName, hashName, cidx,offset,limit,cols...)
}