package routers

import (
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
    "github.com/ouqiang/delay-queue/delayqueue"
)

func Push(resp http.ResponseWriter, req *http.Request)  {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("读取body错误#%s", err.Error())
        return
    }
    var job delayqueue.Job
    err = json.Unmarshal(body, &job)
    if err != nil {
        log.Printf("解析json失败#%s", err.Error())
        return
    }
    delayqueue.Push(job)
}

func Pop(resp http.ResponseWriter, req *http.Request)  {
}

func Delete(resp http.ResponseWriter, req *http.Request)  {
    
}

func Finish(resp http.ResponseWriter, req *http.Request)  {
    
}