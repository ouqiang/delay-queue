package delayqueue

type Job struct  {
    Topic string
    Id string
    Delay int
    TTR int
    Body string
}

// 获取Job
func GetJob(key string) (string, error)  {
    value, err := execRedisCommand("GET", key)
    if err != nil {
        return "", err
    }

    byteValue := value.([]byte)

    return string(byteValue), nil
}

// 添加Job
func PutJob(key, value string) error {
    _, err := execRedisCommand("SET", key, value)

    return err
}

// 删除Job
func RemoveJob(key string) error {
    _, err := execRedisCommand("DEL", key)

    return err
}