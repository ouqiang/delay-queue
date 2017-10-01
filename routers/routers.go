package routers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ouqiang/delay-queue/delayqueue"
)

// TopicRequest Job类型请求json
type TopicRequest struct {
	Topic string `json:"topic"`
}

// IdRequest JobId请求json
type IdRequest struct {
	Id string `json:"id"`
}

// Push 添加job
func Push(resp http.ResponseWriter, req *http.Request) {
	var job delayqueue.Job
	err := readBody(resp, req, &job)
	if err != nil {
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

	if job.Delay <= 0 || job.Delay > (1<<31) {
		resp.Write(generateFailureBody("delay 取值范围1 - (2^31 - 1)"))
		return
	}

	if job.TTR <= 0 || job.TTR > 86400 {
		resp.Write(generateFailureBody("ttr 取值范围1 - 86400"))
		return
	}

	log.Printf("add job#%+v\n", job)
	job.Delay = time.Now().Unix() + job.Delay
	err = delayqueue.Push(job)

	if err != nil {
		log.Printf("添加job失败#%s", err.Error())
		resp.Write(generateFailureBody(err.Error()))
	} else {
		resp.Write(generateSuccessBody("添加成功", nil))
	}
}

// Pop 获取job
func Pop(resp http.ResponseWriter, req *http.Request) {
	var topicRequest TopicRequest
	err := readBody(resp, req, &topicRequest)
	if err != nil {
		return
	}
	topic := strings.TrimSpace(topicRequest.Topic)
	if topic == "" {
		resp.Write(generateFailureBody("topic不能为空"))
		return
	}
	// 多个topic逗号分隔
	topics := strings.Split(topic, ",")
	job, err := delayqueue.Pop(topics)
	if err != nil {
		log.Printf("获取job失败#%s", err.Error())
		resp.Write(generateFailureBody(err.Error()))
		return
	}

	if job == nil {
		resp.Write(generateSuccessBody("操作成功", nil))
		return
	}

	type Data struct {
		Id   string `json:"id"`
		Body string `json:"body"`
	}

	data := Data{
		Id:   job.Id,
		Body: job.Body,
	}

	log.Printf("get job#%+v", data)

	resp.Write(generateSuccessBody("操作成功", data))
}

// Delete 删除job
func Delete(resp http.ResponseWriter, req *http.Request) {
	var idRequest IdRequest
	err := readBody(resp, req, &idRequest)
	if err != nil {
		return
	}
	id := strings.TrimSpace(idRequest.Id)
	if id == "" {
		resp.Write(generateFailureBody("job id不能为空"))
		return
	}

	log.Printf("delete job#jobId-%s\n", id)
	err = delayqueue.Remove(id)
	if err != nil {
		log.Printf("删除job失败#%s", err.Error())
		resp.Write(generateFailureBody(err.Error()))
		return
	}

	resp.Write(generateSuccessBody("操作成功", nil))
}

// Get 查询job
func Get(resp http.ResponseWriter, req *http.Request) {
	var idRequest IdRequest
	err := readBody(resp, req, &idRequest)
	if err != nil {
		return
	}
	id := strings.TrimSpace(idRequest.Id)
	if id == "" {
		resp.Write(generateFailureBody("job id不能为空"))
		return
	}
	job, err := delayqueue.Get(id)
	if err != nil {
		log.Printf("查询job失败#%s", err.Error())
		resp.Write(generateFailureBody(err.Error()))
		return
	}

	if job == nil {
		resp.Write(generateSuccessBody("操作成功", nil))
		return
	}

	resp.Write(generateSuccessBody("操作成功", job))
}

// ResponseBody 响应Body格式
type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func readBody(resp http.ResponseWriter, req *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("读取body错误#%s", err.Error())
		resp.Write(generateFailureBody("读取request body失败"))
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		log.Printf("解析json失败#%s", err.Error())
		resp.Write(generateFailureBody("解析json失败"))
		return err
	}

	return nil
}

func generateSuccessBody(msg string, data interface{}) []byte {
	return generateResponseBody(0, msg, data)
}

func generateFailureBody(msg string) []byte {
	return generateResponseBody(1, msg, nil)
}

func generateResponseBody(code int, msg string, data interface{}) []byte {
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
