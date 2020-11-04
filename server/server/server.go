package server

import (
	"../config"
	"../model"
	"fmt"
	"github.com/kataras/iris"
	"strconv"
)


func Init() error {

	env := config.GetConfig()
	//1、初始化数据库
	err :=  model.InitDB(env.SSDBHost,env.SSDBPort,env.SSDBPwd)
	if err != nil {
		return err
	}
	//2、初始化安网基础设备数据库
	if env.BaseSSDBHost != "" {
		err =  model.InitBaseDeviceSSDB(env.BaseSSDBHost,env.BaseSSDBPort,env.BaseSSDBPwd)
		if err != nil {
			return err
		}
	}

	return nil
}
func Close()  {

	model.CloseDB()
	if config.GetConfig().BaseSSDBHost != ""{
		model.CloseBaseDeviceDB()
	}
}

func Start()  {
	//开启http go服务器
	app := iris.New()
	env := config.GetConfig()

	//服务器请求访问监控入口
	app.Use(func(ctx iris.Context) {
		//打印请求path[method]
		fmt.Println("Path:", ctx.Path(), "[", ctx.Method(), "]")
		ctx.Next()
	})

	//服务器路由 接口设置
	registerInitCgi(app)

	if env.HttpsEnable{
		if env.HttpsPort <= 0 {
			panic("配置文件：https服务器端口设置错误")
		}
		go app.Run(iris.Addr(":" + strconv.Itoa(env.HttpPort)))

		_ = app.Run(iris.TLS(":" + strconv.Itoa(env.HttpsPort), env.HttpsCertPem, env.HttpsCertKey))
	}else{
		if env.HttpPort <= 0{
			panic("配置文件：http服务器端口设置错误")
		}
		_ = app.Run(iris.Addr(":" + strconv.Itoa(env.HttpPort)))
	}
}
