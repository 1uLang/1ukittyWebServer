package server

import (
	"../model"
	"fmt"
	"github.com/kataras/iris"
	"time"
)

//用户 server层
type user struct {
}

func init() {
	initCgiAdd(&user{})
}
func (ptr *user) RelativePath() string {
	return "/user/"
}
func (ptr *user) Init(part iris.Party) {
	//不判断权限和true的区别是,true权限设置之后可以修改用户自身的权限为false来达到黑名单的效果
	//不设置的话则表示所有用户都可以一直访问该接口,且不能通过修改用户自身的权限来阻止其访问
	part.Post("login", ptr.login)
	part.Post("register", ptr.register)

	part.Get("info", ptr.info)
	part.Post("info", ptr.info)
	//默认权限为允许访问,即只有当用户自身权限为false的时候才不能访问(可用于黑名单控制)
	registerDefaultPower(ptr.RelativePath()+"info", "获取用户自己的信息", true)

	part.Post("update", ptr.update)
	part.Get("update", ptr.update)
	registerDefaultPower(ptr.RelativePath()+"update", "修改用户自己的信息", true)
}

func (ptr *user) login(ctx iris.Context) {
	in := struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}{}

	if err := ctx.ReadJSON(&in); err != nil {
		RetError(ctx, err.Error())
		return
	}

	//登录支持
	id, ok, err := model.User.Login(in.Username, in.Password)
	fmt.Println(id, ok, err)
	if err != nil {
		RetError(ctx, err.Error())
		return
	}

	if !ok {
		RetError(ctx, "账号或密码错误")
		return
	}

	tokenString, err := createToken(in.Username, id)
	if err != nil {
		RetError(ctx, "账号或密码错误:"+err.Error())
		return
	}

	RetOk(ctx, struct {
		Expire int64  `json:"expire"`
		Token  string `json:"token"`
	}{
		time.Now().Add(tokenConfig.Expire).Unix(),
		tokenString,
	})

	return
}

//获取自己信息
func (ptr *user) info(ctx iris.Context) {
	user, _ := GetUserId(ctx)
	ret, err := model.User.Info(user)
	if err != nil {
		RetError(ctx, err.Error())
	} else {
		RetOk(ctx, ret)
	}
}
//新增用户
func (ptr *user) register(ctx iris.Context) {
	in := model.User.GetShowInfo()
	err := ctx.ReadJSON(&in)
	if err != nil {
		RetError(ctx, err.Error())
		return
	}

	if in.Account == "" || in.Password == ""{
		RetError(ctx, "参数错误")
		return
	}

	err = model.User.Add(in)
	if err != nil {
		RetError(ctx, err.Error())
	} else {
		RetOk(ctx, nil)
	}
}

//修改用户自己的信息
func (ptr *user) update(ctx iris.Context) {
	in := model.User.GetShowInfo()
	err := ctx.ReadJSON(&in)
	if err != nil {
		RetError(ctx, err.Error())
		return
	}

	if in.ID == "" {
		RetError(ctx, "Id错误")
		return
	}

	user, _ := GetUserId(ctx)
	err = model.User.Update(user, in)
	if err != nil {
		RetError(ctx, err.Error())
	} else {
		RetOk(ctx, nil)
	}
}
