package main

import (
	"cmdtools/cmdhandler"
	"cmdtools/core/dragon"
	"os"
)

func main() {
	// 读取配置文件
	dragon.AppInit()

	// 判断传入参数格式是否正确
	cmdhandler.HandleArgs(os.Args)

}
