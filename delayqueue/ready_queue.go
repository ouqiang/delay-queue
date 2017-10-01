package delayqueue

import (
	"fmt"

	"github.com/ouqiang/delay-queue/config"
)

// 添加JobId到队列中
func pushToReadyQueue(queueName string, jobId string) error {
	queueName = fmt.Sprintf(config.Setting.QueueName, queueName)
	_, err := execRedisCommand("RPUSH", queueName, jobId)

	return err
}

// 从队列中阻塞获取JobId
func blockPopFromReadyQueue(queues []string, timeout int) (string, error) {
	var args []interface{}
	for _, queue := range queues {
		queue = fmt.Sprintf(config.Setting.QueueName, queue)
		args = append(args, queue)
	}
	args = append(args, timeout)
	value, err := execRedisCommand("BLPOP", args...)
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
