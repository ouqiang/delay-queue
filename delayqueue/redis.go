package delayqueue

import (
    "time"
    "github.com/ouqiang/delay-queue/config"
    "github.com/garyburd/redigo/redis"
    "log"
)

var (
    RedisPool *redis.Pool
)

// 初始化连接池
func initRedisPool() *redis.Pool {
    pool := &redis.Pool{
        MaxIdle: config.Setting.Redis.MaxIdle,
        MaxActive: config.Setting.Redis.MaxActive,
        IdleTimeout: 300 * time.Second,
        Dial: redisDial,
        TestOnBorrow: redisTestOnBorrow,
        Wait: true,
    }

    return pool
}

// 连接redis
func redisDial() (redis.Conn, error)  {
    conn, err := redis.Dial("tcp", config.Setting.Redis.Host)
    if err != nil {
        log.Printf("连接redis失败#%s", err.Error())
        return nil, err
    }

    if (config.Setting.Redis.Password == "") {
       return conn, err
    }

    if _, err := conn.Do("AUTH", config.Setting.Redis.Password); err != nil {
        conn.Close()
        log.Printf("redis认证失败#%s", err.Error())
        return nil, err
    }

    return conn, err
}

// 从池中取出连接后，判断连接是否有效
func redisTestOnBorrow(conn redis.Conn, t time.Time) error {
    if time.Since(t) < time.Minute {
        return nil
    }
    _, err := conn.Do("PING")

    return err
}

// 执行redis命令, 执行完成后连接自动放回连接池
func execRedisCommand(command string, args ...interface{}) (interface{}, error) {
    redis := RedisPool.Get()
    defer redis.Close()

    return redis.Do(command, args...)
}
