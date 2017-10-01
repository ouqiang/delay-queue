package delayqueue

import (
	"strconv"
)

// BucketItem bucket中的元素
type BucketItem struct {
	timestamp int64
	jobId     string
}

// 添加JobId到bucket中
func pushToBucket(key string, timestamp int64, jobId string) error {
	_, err := execRedisCommand("ZADD", key, timestamp, jobId)

	return err
}

// 从bucket中获取延迟时间最小的JobId
func getFromBucket(key string) (*BucketItem, error) {
	value, err := execRedisCommand("ZRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	var valueBytes []interface{}
	valueBytes = value.([]interface{})
	if len(valueBytes) == 0 {
		return nil, nil
	}
	timestampStr := string(valueBytes[1].([]byte))
	item := &BucketItem{}
	item.timestamp, _ = strconv.ParseInt(timestampStr, 10, 64)
	item.jobId = string(valueBytes[0].([]byte))

	return item, nil
}

// 从bucket中删除JobId
func removeFromBucket(bucket string, jobId string) error {
	_, err := execRedisCommand("ZREM", bucket, jobId)

	return err
}
