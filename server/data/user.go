package data

import (
	"../tools"
	"fmt"
	"github.com/seefan/gossdb"
	"github.com/seefan/sorm"
	"os"
)

type UserInfo struct {
	ID       string `json:"id"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Tel      string `json:"tel"`
}
type userCols struct {
	Account  string `tbl:"!DATA"`
	Password string `tbl:"!DATA"`
	Name     string `tbl:"DATA"`
	Tel      string `tbl:"DATA"`
}
type userCidx struct {
	Account []string `cidx:"only"`
}

func (c UserInfo) TblName() string {
	return "secnetServer_userData"
}

func (c UserInfo) TblType() string {
	return "INCR"
}

//set ssdb cols
var usercols = &userCols{}
var usercidx = &userCidx{}

func (c UserInfo) Cols() *userCols {
	return usercols
}
func (c UserInfo) Cidx() *userCidx {
	return usercidx
}

/*------------------------------------- private -----------------------------------------*/
func init() {
	err := sorm.Sync(&UserInfo{})
	if err != nil {
		fmt.Printf("sorm sync UserInfo:%s\n", err.Error())
		os.Exit(1)
	}
}
func UserInfoCreateTbl(pool *gossdb.Connectors) {
	db, err := pool.NewClient()
	if err != nil {
		fmt.Printf("UserInfoCreateTbl pool new client err:%s\n", err)
	}
	defer db.Close()
	err = sorm.CreateTables(db, (&UserInfo{}).TblName())
	if err != nil {
		fmt.Printf(" UserInfoCreateTbl UserInfo tbl error:%s\n", err)
	}
}
func (c UserInfo) IsExistAccount(db *gossdb.Client) (bool, error) {

	if c.Account == "" {
		return false, fmt.Errorf("用户账户能为空")
	}
	cidx := gossdb.Cidx_scan_t{c.Cidx().Account, 0, nil, c.Account, ""}

	_, ret, err := db.Tcidx_scan(c.TblName(), cidx, 0, 1)
	if err != nil {
		return false, err
	}

	return len(ret) > 0 && ret[0].Cidx[c.Cols().Account] == c.Account, nil
}
func (c UserInfo) IsExist(db *gossdb.Client) (bool, error) {

	if c.ID != "" {
		val, err := db.Tcol_count(c.TblName(), c.ID)
		if err != nil {
			return false, err
		}
		return val > 0, err
	}
	return c.IsExistAccount(db)
}
func (c *UserInfo) Add(db *gossdb.Client) error {

	if c.Account == "" || c.Password == "" {
		return fmt.Errorf("用户账号或密码不能为空")
	}
	ok, err := c.IsExist(db)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("账号已注册")
	}
	orm := sorm.Engine{}
	orm.DB = db
	orm.SetTable(sorm.Tables[c.TblName()])
	id, err := orm.Tinsert(c)
	if err != nil {
		return err
	}
	c.ID = id[0]
	return err
}

func (c *UserInfo) Get(db *gossdb.Client) error {

	if c.ID == "" {
		return fmt.Errorf("用户id不能为空")
	}
	val, err := db.Tcol_count(c.TblName(), c.ID)
	if err != nil {
		return err
	}
	if val <= 0 {
		return fmt.Errorf("无效用户id，该用户不存在")
	}
	vals, err := db.Multi_tcol_get(c.TblName(), c.ID, c.Cols().Account, c.Cols().Name, c.Cols().Tel)
	if err != nil {
		return err
	}

	c.Account = vals[0].String()
	c.Name = vals[1].String()
	c.Tel = vals[2].String()

	return nil
}
func (c UserInfo) Update(db *gossdb.Client) error {
	if c.ID == "" {
		return fmt.Errorf("用户id不能为空")
	}
	ok, err := c.IsExist(db)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("无效用户id，该用户id不存在")
	}
	//账号不能修改
	c.Account = ""

	orm := sorm.Engine{}
	orm.DB = db
	orm.SetTable(sorm.Tables[c.TblName()])
	return orm.Tupdate(c.ID, c)
}
func (c UserInfo) SetPassword(db *gossdb.Client) error {
	if c.ID == "" {
		return fmt.Errorf("用户id不能为空")
	}
	if c.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return db.Tupdate(c.TblName(), c.ID, map[string]interface{}{c.Cols().Password: c.Password})
}
func (c UserInfo) Del(db *gossdb.Client) error {
	if c.ID == "" {
		return fmt.Errorf("用户id不能为空")
	}
	ok, err := c.IsExist(db)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("无效用户id，该用户不存在")
	}

	_, err = db.Trow_clear(c.TblName(), c.ID)
	return err
}

func (c UserInfo) Login(db *gossdb.Client) (string, bool, error) {
	if c.Password == "" || c.Account == "" {
		return "", false, fmt.Errorf("用户账户密码不能为空")
	}
	cidx := gossdb.Cidx_scan_t{c.Cidx().Account, 0, nil, c.Account, ""}

	_, ret, err := db.Tcidx_scan(c.TblName(), cidx, 0, 1, c.Cols().Password)
	if err != nil {
		return "", false, err
	}
	if len(ret) == 0 || ret[0].Cidx[c.Cols().Account] != c.Account {
		return "", false, nil
	}
	return ret[0].Idx, ret[0].Colv[c.Cols().Password].String() == c.Password, nil
}
func (c UserInfo) List(db *gossdb.Client, limit, page int64) (num int64, list []UserInfo, err error) {

	num, err = db.Trow_count(c.TblName())
	if err != nil {
		return
	}
	cidx := gossdb.Cidx_scan_t{c.Cidx().Account, 0, nil, c.Account, tools.StringAdd(c.Account)}

	offset := limit * (page - 1)
	if offset > num {
		offset = 0
	}
	count, ret, err := db.Tcidx_scan(c.TblName(), cidx, offset, limit, "*")

	if err != nil {
		return
	}
	list = make([]UserInfo, count)
	for k, v := range ret {
		list[k].ID = v.Idx
		list[k].Account = v.Cidx[c.Cols().Account]
		//list[k].Password = v.Colv[c.Cols().Password].String()

		list[k].Name = v.Colv[c.Cols().Name].String()
		list[k].Tel = v.Colv[c.Cols().Tel].String()
	}
	return
}

func (c UserInfo) FindByAccount(db *gossdb.Client) (string, error) {

	cidx := gossdb.Cidx_scan_t{c.Cidx().Account, 0, nil, c.Account, ""}

	_, ret, err := db.Tcidx_scan(c.TblName(), cidx, 0, 1)

	if err != nil {
		return "", err
	}
	if len(ret) > 0 && ret[0].Cidx[c.Cols().Account] == c.Account {
		return ret[0].Idx, nil
	}
	return "", nil
}

func (c UserInfo) GetAccount(db *gossdb.Client) (string, error) {
	if c.ID == "" {
		return "", fmt.Errorf("用户id不能为空")
	}
	val, err := db.Tcol_count(c.TblName(), c.ID)
	if err != nil {
		return "", err
	}
	if val <= 0 {
		return "", fmt.Errorf("无效用户id，该用户不存在")
	}

	ret, err := db.Tcol_get(c.TblName(), c.ID, c.Cols().Account)
	return ret.String(), err
}
