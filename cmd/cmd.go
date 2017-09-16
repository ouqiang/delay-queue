package cmd

import (
	"flag"
	"fmt"
	"github.com/ouqiang/delay-queue/config"
	"github.com/ouqiang/delay-queue/delayqueue"
	"github.com/ouqiang/delay-queue/routers"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type Cmd struct{}

var (
	version    bool
	configFile string
)

const (
	AppVersion = "0.3"
)

func (cmd *Cmd) Run() {

	// 解析命令行参数
	cmd.parseCommandArgs()
	if version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	// 初始化配置
	config.Init(configFile)
	// 初始化队列
	delayqueue.Init()

	// 运行web server
	cmd.runWeb()
}

func (cmd *Cmd) parseCommandArgs() {
	// 配置文件
	flag.StringVar(&configFile, "c", "", "./delay-queue -c /path/to/delay-queue.conf")
	// 版本
	flag.BoolVar(&version, "v", false, "./delay-queue -v")
	flag.Parse()
}

func (cmd *Cmd) runWeb() {
	http.HandleFunc("/push", routers.Push)
	http.HandleFunc("/pop", routers.Pop)
	http.HandleFunc("/finish", routers.Delete)
	http.HandleFunc("/delete", routers.Delete)

	log.Printf("listen %s\n", config.Setting.BindAddress)
	err := http.ListenAndServe(config.Setting.BindAddress, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
