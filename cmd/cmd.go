package cmd

import (
    "flag"
    "fmt"
    "os"
    "net/http"
    "github.com/ouqiang/delay-queue/config"
    "github.com/ouqiang/delay-queue/routers"
    "github.com/ouqiang/delay-queue/dealyqueue"
    "syscall"
    "os/signal"
    "log"
)

type Cmd struct {}

var (
    version bool
    configFile string
)

const (
    AppVersion = "0.1"
)

func (cmd *Cmd) Run()  {
    // 解析命令行参数
    cmd.parseCommandArgs();
    if version {
        fmt.Println(AppVersion)
        os.Exit(0)
    }
    // 初始化配置
    config.Init(configFile)
    // 初始化队列
    delayqueue.Init()

    // 捕捉信号
    go catchSignal();

    // 运行web server
    cmd.runWeb()
}

func (cmd *Cmd) parseCommandArgs()  {
    // 配置文件
    flag.StringVar(&configFile, "c", "", "./delay-queue -c /path/to/delay-queue.conf")
    // 版本
    flag.BoolVar(&version, "v", false, "./delay-queue -v")
    flag.Parse()
}

func (cmd *Cmd) runWeb()  {
    http.HandleFunc("/push", routers.Push)
    http.HandleFunc("/pop", routers.Pop)
    http.HandleFunc("/finish", routers.Finish)
    http.HandleFunc("/delete", routers.Delete)

    http.ListenAndServe(config.Setting.BindAddress, nil)
}

// 捕捉信号
func catchSignal()  {
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    for {
        s := <- c
        switch s {
        case syscall.SIGHUP:
            log.Println("收到SIGNUP信号, 忽略")
        case  syscall.SIGINT, syscall.SIGTERM:
            shutdown()
        }
    }
}

func shutdown()  {
    // 释放连接池资源
    delayqueue.RedisPool.Close()

    os.Exit(0)
}