package delayqueue

import (
    "time"
    "github.com/ouqiang/delay-queue/config"
    "fmt"
    "log"
    "strings"
)

var (
    // 每个定时器对应一个bucket
    timers []*time.Ticker
    // 保存待放入bucket中的job
    queue chan Job
)

func Init()  {
    RedisPool = initRedisPool()
    initTimers()
    queue = make(chan Job, 100)
    go waitJob()
}

// 添加一个Job到队列中
func Push(job Job)  {
    job.Id = strings.TrimSpace(job.Id)
    job.Topic = strings.TrimSpace(job.Topic)
    job.Body = strings.TrimSpace(job.Body)

    if job.Id == "" || job.Topic == "" || job.Delay < 0 || job.TTR <= 0 {
        return
    }

    queue <- job
}

func waitJob()  {
    bucketNameChan := generateBucketName()
    for job := range queue {
        putJob(job.Id, job)
        delayTimestamp := time.Now().Unix() + job.Delay
        pushToBucket(<-bucketNameChan, delayTimestamp, job.Id)
    }
}

// 轮询获取Job名称, 使job分布到不同bucket中, 提高扫描速度
func generateBucketName() (chan string) {
    c := make(chan string)
    go func() {
        i := 1
        for {
            c <- fmt.Sprintf("dq_bucket_%d", i)
            if i >= config.Setting.Bucket {
                i = 1
            } else {
                i++
            }
        }
    }()

    return c
}

// 初始化定时器
func initTimers()  {
    timers = make([]*time.Ticker, config.Setting.Bucket)
    var bucketName string
    for i := 0; i < config.Setting.Bucket; i++ {
        timers[i] = time.NewTicker(1 * time.Second)
        bucketName = fmt.Sprintf("dq_bucket_%d", i + 1)
        go waitTicker(timers[i], bucketName)
    }
}

func waitTicker(timer *time.Ticker, bucketName string)  {
    for {
        select {
            case t := <- timer.C:
            tickHandler(t, bucketName)
        }
    }
}

// 扫描bucket, 取出延迟时间小于当前时间的Job
func tickHandler(t time.Time, bucketName string)  {
    for {
        bucketItem, err := getFromBucket(bucketName)
        if err != nil {
            log.Printf("扫描bucket错误#bucket-%s#%s", bucketName, err.Error())
            return
        }

        // 集合为空
        if bucketItem == nil {
            return
        }

        // 延迟时间未到
        if bucketItem.timestamp > t.Unix() {
            return
        }

        // 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
        job, err := getJob(bucketItem.jobId)
        if err != nil {
            log.Printf("获取Job元信息失败#bucket-%s#%s", bucketName, err.Error())
        }

        err = pushToReadyQueue(job.Topic, job.Id)
        if err != nil {
            log.Printf("JobId放入ready queue失败#bucket-%s#job-%+v#%s",
                bucketName, job, err.Error())
            return
        }

        // 从bucket中删除

        removeFromBucket(bucketName, job.Id)
    }
}
