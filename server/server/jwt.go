package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"time"
)

type tokenArgs struct {
	Secret string
	Expire time.Duration
}
var tokenConfig = tokenArgs{
	"7sdml7VjNfJFL5tvup7hK0KJcYQUGV",
	12 * time.Hour * time.Duration(1),
}
var jwtHandler = jwtMiddleware.New(jwtMiddleware.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenConfig.Secret), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(ctx iris.Context, err error) {
		fmt.Printf("{%s}访问{%s}时检查token失败:%s\n", ctx.RemoteAddr(), ctx.Path(), err)
		RetUnauthorized(ctx)
		return
	},
	Expiration: true,
})


func createToken(username, id string) (tokenstr string, err error) {
	//生成加密串过程
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nick_name": "iris",
		"email":     "qhe@secnet.cn",
		"id":        id,
		"name":      username,
		"iss":       "Iris",
		"iat":       time.Now().Unix(),
		"jti":       "qian",
		"exp":       time.Now().Add(tokenConfig.Expire).Unix(),
	})
	//把token已约定的加密方式和加密秘钥加密，当然也可以使用不对称加密
	return token.SignedString([]byte(tokenConfig.Secret))
	//登录时候，把tokenString返回给客户端，然后需要登录的页面就在header上面附此字符串
	//eg: header["Authorization"] = "bears "+tokenString
}
//必须校验过token才能调用这个函数
func GetTokenInfo(ctx iris.Context) (info map[string]interface{}) {
	//token := ctx.Values().Get("jwt").(*jwt.Token)
	info = ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	//info["id"].(float64)
	//info["nick_name"].(string)
	return
}

func GetUserId(ctx iris.Context) (id string, isExist bool) {

	//return "root",true
	id, isExist = GetTokenInfo(ctx)["id"].(string)

	return
}
func GetUserName(ctx iris.Context) (name string, isExist bool) {

	//return "root",true
	name, isExist = GetTokenInfo(ctx)["name"].(string)

	return
}