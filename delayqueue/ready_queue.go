package delayqueue

import (
    "fmt"
    "github.com/ouqiang/delay-queue/config"
)

type ReadyQueue struct {}

func pushToReadyQueue(queueName string, jobId string) error {
    queueName = fmt.Sprintf(config.Setting.QueueName, queueName)
    _, err := execRedisCommand("RPUSH", queueName, jobId)

    return err
}

func blockPopFromReadyQueue(queueName string, timeout int) (string, error){
    queueName = fmt.Sprintf(config.Setting.QueueName, queueName)
    value, err := execRedisCommand("BLPOP", queueName, timeout)
    if err != nil {
        return "", err
    }
    if value == nil {
        return "", nil
    }
    var valueBytes []interface{}
    valueBytes = value.([]interface{})
    if len(valueBytes) == 0 {
        return "", nil
    }
    element := string(valueBytes[1].([]byte))

    return element, nil
}
