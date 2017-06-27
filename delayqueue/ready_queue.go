package delayqueue


type ReadyQueue struct {}

func pushToReadyQueue(queueName string, jobId string) error {
    _, err := execRedisCommand("RPUSH", queueName, jobId)

    return err
}

func popFromReadyQueue(queueName string) (string, error){
    value, err := execRedisCommand("LPOP", queueName)
    if err != nil {
        return "", err
    }
    if value == nil {
        return "", nil
    }
    byteValue := value.([]byte)

    return string(byteValue), nil
}
