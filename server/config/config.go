package config

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose/utils"
	"io/ioutil"
	"os"
)

type config struct {
	Url           	string `json:"url"`
	Host          	string `json:"host"`
	HttpPort      	int    `json:"httpPort"`
	HttpsPort     	int    `json:"httpsPort"`

	HttpsEnable 	bool 	`json:"https"`
	HttpsCertPem	string 	`json:"https_cert_pem"`
	HttpsCertKey	string	`json:"https_cert_key"`

	SSDBHost 		string 	`json:"SSDBHost"`
	SSDBPort 		int 	`json:"SSDBPort"`
	SSDBPwd 		string 	`json:"SSDBPwd"`

	BaseSSDBHost 		string 	`json:"baseSSDBHost"`
	BaseSSDBPort 		int 	`json:"baseSSDBPort"`
	BaseSSDBPwd 		string 	`json:"baseSSDBPwd"`
}
var env =  config{}

func GetConfig() config{
	return env
}
func ReadFile(filename string)error  {

	if filename == "" {
		return fmt.Errorf("配置文件不能为空")
	}
	fmt.Println("read config file : ", filename)

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &env)
	if err != nil {
		return err
	}
	if env.SSDBHost == "" {
		panic("配置文件错误 未配置SSDB数据库")
	}
	if env.HttpsPort  > 0 {
		env.HttpsEnable = true
		if env.HttpsCertKey != "" {
			IsExist := utils.FileExists(env.HttpsCertKey)
			if !IsExist {
				return fmt.Errorf("https CERT Key文件[%s] 不存在",env.HttpsCertKey)
			}
		}
		if env.HttpsCertPem != "" {
			IsExist := utils.FileExists(env.HttpsCertPem)
			if !IsExist {
				return fmt.Errorf("https CERT Pem文件[%s] 不存在",env.HttpsCertPem)
			}
		}
	}
	return nil
}