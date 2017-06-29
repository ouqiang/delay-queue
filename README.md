# delay-queue
延迟队列, 参考[有赞延迟队列设计](http://tech.youzan.com/queuing_delay)实现


## 下载
* [Linux-64位](http://opns468ov.bkt.clouddn.com/delay-queue-linux-amd64.tar.gz)
* [Mac OS-64位](http://opns468ov.bkt.clouddn.com/delay-queue-darwin-amd64.tar.gz)
* [Windows-64位](http://opns468ov.bkt.clouddn.com/delay-queue-windows-amd64.zip)


## 源码安装
* `go`语言版本1.7+
* `go get -d github.com/ouqiang/delay-queue`
* `go build`


## 运行
`./delay-queue -c delay-queue.conf`, 默认监听 `0.0.0.0:9277`

## HTTP接口

* 请求方法 `POST`   
* 请求Body及返回值均为`json`

### 返回值
```json
{
  "code": 0,
  "message": "添加成功",
  "data": null
}
```

|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   code  |     int     |    状态码    | 0: 成功 非0: 失败 |
| message |    string   | 状态描述信息 |                   |
|   data  | object, null |   附加信息   |                   |

### 添加任务   
URL地址 `/push`   
```json
{
  "topic": "order",
  "id": "15702398321",
  "delay": 3600,
  "ttr": 120,
  "body": "{\"uid\": 10829378,\"created\": 1498657365 }"
}
```
|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   topic  | string     |    Job类型                   |                     |
|   id     | string     |    Job唯一标识                   |                   |
|   delay  | int        |    Job需要延迟的时间, 单位：秒    |                   |
|   ttr  | int        |    Job执行超时时间, 单位：秒   |                   |
|   body   | string     |    Job的内容，供消费者做具体的业务处理，如果是json格式需转义 |                   |

### 获取ready queue中的任务    
```json
{
  "topic": "order"
}
```
|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   topic  | string     |    Job类型                   |                     |


队列中有任务返回值
```json
{
  "code": 0,
  "message": "操作成功",
  "data": {
    "id": "15702398321",
    "body": "{\"uid\": 10829378,\"created\": 1498657365 }"
  }
}
```
队列为空返回值   
```json
{
  "code": 0,
  "message": "操作成功",
  "data": null
}
```


### 删除任务  
URL地址 `/delete`   

```json
{
  "id": "15702398321"
}
```

|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   id  | string     |    Job唯一标识       |            |

  
### 完成任务   
URL地址 `/finish`   

```json
{
  "id": "15702398321"
}
```

|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   id  | string     |    Job唯一标识    |                     |

