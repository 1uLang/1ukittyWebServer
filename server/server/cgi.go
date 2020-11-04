package server

import (
	"fmt"
	"github.com/kataras/iris"
	"strings"
	"sync"
)

var cgi = struct {
	WhiteList       map[string]bool //精确匹配
	WhiteListPrefix map[string]bool //匹配前缀
}{
	WhiteList: map[string]bool{
		"/login":               true,
		"/user/login":          true,
		"/user/register":          true,
	},
	WhiteListPrefix: map[string]bool{
		"/file/":       true,
	},
}

func checkWhiteList(path string) (isWhite bool) {
	isWhite = false
	if path == "" {
		return
	}

	_, exist := cgi.WhiteList[path]
	if exist {
		isWhite = true
		return
	}

	for k := range cgi.WhiteListPrefix {
		if strings.HasPrefix(path, k) {
			isWhite = true
			fmt.Println(path + "前缀白名单请求")
			return
		}
	}
	return
}

func before(ctx iris.Context) {

	fmt.Println(ctx.Path(), ctx.Method())
	//跨域
	ctx.Header("Access-Control-Allow-Origin","*")
	ctx.Header("Access-Control-Allow-Credentials","true")
	ctx.Header("Access-Control-Allow-Headers","Access-Control-Allow-Origin,Content-Type")

	//放行options
	if ctx.Method() == "OPTIONS" {
		ctx.Next()
		return
	}

	////检查白名单
	if checkWhiteList(ctx.Path()) {
		ctx.Next()
		return
	}

	//检查token是否合法
	if err := jwtHandler.CheckJWT(ctx); err != nil {
		jwtHandler.Config.ErrorHandler(ctx, err)
		return
	}

	//判断后台白名单
	_, isExist := GetUserName(ctx)
	if !isExist {
		panic("token 有逻辑错误1")
	}
	ok, err := CheckPower(ctx.Path())
	if err != nil {
		RetError(ctx, "检查权限时:"+err.Error())
		return
	}
	if !ok {
		fmt.Printf("{ip:%s}{method:%v} 访问{%s}:无权限拒绝访问\n",
			ctx.RemoteAddr(), ctx.Method(), ctx.Path())
		RetUnPower(ctx)
		return
	}
	ctx.Next()
}

type CgiRetInfo struct {
	Code         int         `json:"code"`
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage string      `json:"message,omitempty"`
}

type CgiRetData struct {
	DataList  interface{} `json:"DataList"`
	Count     int64       `json:"Count"`
	PageIndex int64       `json:"PageIndex"`
	PageSize  int64       `json:"PageSize"`

	TotalPage int64 `json:"TotalPage"`
	CurrPage  int64 `json:"CurrPage"`
}

type cgiInterface interface {
	Init(part iris.Party)
	RelativePath() string
}

var initCgiList []cgiInterface
var initCgiMutex sync.Mutex

//注册init cgi
func initCgiAdd(cgi cgiInterface) {
	initCgiMutex.Lock()
	defer initCgiMutex.Unlock()
	initCgiList = append(initCgiList, cgi)
}
func registerInitCgi(app *iris.Application) {
	//所有请求执行的第一个处理程序
	app.Use(before)

	initCgiMutex.Lock()
	for _, cgi := range initCgiList {
		fmt.Println("init cgi:", cgi.RelativePath())
		cgi.Init(app.Party(cgi.RelativePath()))
	}
	initCgiMutex.Unlock()

	//其它请求
	app.Any("{path:path}", RetNotFount)
}

func RetOk(ctx iris.Context, data interface{}, msg ...string) {
	if len(msg) > 0 {
		_, _ = ctx.JSON(CgiRetInfo{Code: 200, Data: data, ErrorMessage: msg[0]})
	} else {
		_, _ = ctx.JSON(CgiRetInfo{Code: 200, Data: data})
	}
}

func RetNotFount(ctx iris.Context) {
	fmt.Printf("{ip:%s}{method:%v}访问{%s}404:未找到资源\n", ctx.RemoteAddr(), ctx.Method(), ctx.Path())
	_, _ = ctx.JSON(CgiRetInfo{Code: 404, ErrorMessage: "资源未找到"})
}

func RetUnauthorized(ctx iris.Context) {
	_, _ = ctx.JSON(CgiRetInfo{Code: 401, ErrorMessage: "登录失效，请重新登录。"})
}
func RetUnPower(ctx iris.Context) {
	_, _ = ctx.JSON(CgiRetInfo{Code: 401, ErrorMessage: "无权限"})
}

func RetError(ctx iris.Context, err string) {
	_, _ = ctx.JSON(CgiRetInfo{Code: 500, ErrorMessage: err})
}
