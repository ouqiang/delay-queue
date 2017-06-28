package routers

import (
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
    "github.com/ouqiang/delay-queue/delayqueue"
    "time"
    "strings"
)

// 添加job
func Push(resp http.ResponseWriter, req *http.Request)  {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("读取body错误#%s", err.Error())
        resp.Write(generateFailureBody("读取request body失败"))
        return
    }
    var job delayqueue.Job
    err = json.Unmarshal(body, &job)
    if err != nil {
        log.Printf("解析json失败#%s", err.Error())
        resp.Write(generateFailureBody("解析json失败"))
        return
    }

    job.Id = strings.TrimSpace(job.Id)
    job.Topic = strings.TrimSpace(job.Topic)
    job.Body = strings.TrimSpace(job.Body)
    if job.Id == "" {
        resp.Write(generateFailureBody("job id不能为空"))
        return
    }
    if job.Topic == "" {
        resp.Write(generateFailureBody("topic 不能为空"))
        return
    }

    if job.Delay <= 0 || job.Delay > (1 << 31) {
        resp.Write(generateFailureBody("delay 取值范围1 - (2^31 - 1)"))
        return
    }

    if job.TTR <= 0 || job.TTR > 86400 {
        resp.Write(generateFailureBody("ttr 取值范围1 - 86400"))
        return
    }

    log.Printf("add job#%+v\n", job)
    job.Delay = time.Now().Unix() + job.Delay
    delayqueue.Push(job)

    resp.Write(generateSuccessBody("添加成功", nil))
}

// 获取job
func Pop(resp http.ResponseWriter, req *http.Request)  {
    topic := req.PostFormValue("topic")
    topic = strings.TrimSpace(topic)
    if topic == "" {
        resp.Write(generateFailureBody("topic不能为空"))
        return
    }
    job, err := delayqueue.Pop(topic)
    if err != nil {
        log.Printf("获取job失败#%s", err.Error())
        resp.Write(generateFailureBody("获取失败"))
        return
    }
    if job == nil {
        resp.Write(generateSuccessBody("操作成功", nil))
        return
    }

    type Data struct  {
        Id string `json:"id"`
        Body string `json:"body"`
    }

    data := Data{
        Id: job.Id,
        Body: job.Body,
    }

    log.Printf("获取job#%+v", data)

    resp.Write(generateSuccessBody("操作成功", data))
}

// 删除job
func Delete(resp http.ResponseWriter, req *http.Request)  {
    id := req.PostFormValue("id")
    id = strings.TrimSpace(id)
    if id == "" {
        resp.Write(generateFailureBody("job id不能为空"))
        return
    }

    err := delayqueue.Remove(id)
    if err != nil {
        resp.Write(generateFailureBody("删除失败"))
        return
    }
    log.Printf("delete job#jobId-%s\n", id)

    resp.Write(generateSuccessBody("操作成功", nil))
}

type ResponseBody struct {
    Code int       `json:"code"`
    Message string `json:"message"`
    Data interface{} `json:"data"`
}

func generateSuccessBody(msg string, data interface{}) ([]byte)  {
     return generateResponseBody(0, msg, data)
}

func generateFailureBody(msg string) ([]byte) {
    return generateResponseBody(1, msg, nil)
}

func generateResponseBody(code int, msg string, data interface{}) ([]byte)  {
    body := &ResponseBody{}
    body.Code = code
    body.Message = msg
    body.Data = data

    bytes, err := json.Marshal(body)
    if err != nil {
        log.Printf("生成response body,转换json失败#%s", err.Error())
        return []byte(`{"code":"1", "message": "生成响应body异常", "data":[]}`)
    }

    return bytes
}