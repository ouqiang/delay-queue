package delayqueue

type ReadyQueue struct {}

func PushToReadyQueue(queueName string, value string) error {
    _, err := execRedisCommand("RPUSH", queueName, value)

    return err
}

func PopFromReadyQueue(queueName string) (string, error){
    value, err := execRedisCommand("LPOP", queueName)
    if err != nil {
        return "", err
    }
    byteValue := value.([]byte)

    return string(byteValue), nil
}
