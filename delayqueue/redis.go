package delayqueue

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/ouqiang/delay-queue/config"
)

var (
	RedisPool *redis.Pool
)

// 初始化连接池
func initRedisPool() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:      config.Setting.Redis.MaxIdle,
		MaxActive:    config.Setting.Redis.MaxActive,
		IdleTimeout:  300 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisTestOnBorrow,
		Wait:         true,
	}

	return pool
}

// 连接redis
func redisDial() (redis.Conn, error) {
	conn, err := redis.Dial(
		"tcp",
		config.Setting.Redis.Host,
		redis.DialConnectTimeout(time.Duration(config.Setting.Redis.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(config.Setting.Redis.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(config.Setting.Redis.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		log.Printf("连接redis失败#%s", err.Error())
		return nil, err
	}

	if config.Setting.Redis.Password != "" {
		if _, err := conn.Do("AUTH", config.Setting.Redis.Password); err != nil {
			conn.Close()
			log.Printf("redis认证失败#%s", err.Error())
			return nil, err
		}
	}

	_, err = conn.Do("SELECT", config.Setting.Redis.Db)
	if err != nil {
		conn.Close()
		log.Printf("redis选择数据库失败#%s", err.Error())
		return nil, err
	}

	return conn, nil
}

// 从池中取出连接后，判断连接是否有效
func redisTestOnBorrow(conn redis.Conn, t time.Time) error {
	_, err := conn.Do("PING")
	if err != nil {
		log.Printf("从redis连接池取出的连接无效#%s", err.Error())
	}

	return err
}

// 执行redis命令, 执行完成后连接自动放回连接池
func execRedisCommand(command string, args ...interface{}) (interface{}, error) {
	redis := RedisPool.Get()
	defer redis.Close()

	return redis.Do(command, args...)
}
