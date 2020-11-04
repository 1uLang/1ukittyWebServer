package main

import (
	"./config"
	"./server"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var godaemon = flag.Bool("d", false, "run app as a daemon with -d=true")

var configFile = flag.String("f", "", "set config file with -f='config file path'")

var ConfigFile = ""

const configPath = "/go/config/config.json"
func init() {

	ConfigFile,_ = os.Getwd()
	fmt.Println("当前工作目录：",ConfigFile)
	ConfigFile += configPath
	if !flag.Parsed() {
		flag.Parse()
	}
	if *configFile != "" {
		ConfigFile = *configFile
	}
	if *godaemon {
		args := os.Args[1:]
		i := 0
		for ; i < len(args); i++ {
			if args[i] == "-d=true" {
				args[i] = "-d=false"
				break
			}
		}
		cmd := exec.Command(os.Args[0], args...)
		cmd.Start()
		fmt.Println("[PID]", cmd.Process.Pid)
		os.Exit(0)
	}
}

func main() {

	//加载配置文件
	if err := config.ReadFile(ConfigFile);err != nil{
		panic(err)
	}
	//初始化服务器
	if err := server.Init();err != nil{
		panic(err)
	}
	defer server.Close()

	//启动 服务器
	server.Start()
}
