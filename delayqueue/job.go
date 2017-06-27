package delayqueue

import "encoding/json"

type Job struct  {
    Topic string `json:"topic"`
    Id string `json:"id"`
    Delay int64 `json:"delay"`
    TTR int `json:"ttr"`
    Body string `json:"body"`
}

// 获取Job
func getJob(key string) (*Job, error)  {
    value, err := execRedisCommand("GET", key)
    if err != nil {
        return nil, err
    }
    if value == nil {
        return nil, nil
    }

    byteValue := value.([]byte)
    job := &Job{}
    err = json.Unmarshal(byteValue, job)
    if err != nil {
        return nil, err
    }

    return job, nil
}

// 添加Job
func putJob(key string, job Job) error {
    value, err := json.Marshal(job)
    if err != nil {
        return err
    }
    _, err = execRedisCommand("SET", key, value)

    return err
}

// 删除Job
func removeJob(key string) error {
    _, err := execRedisCommand("DEL", key)

    return err
}