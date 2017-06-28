package config

import (
    "gopkg.in/ini.v1"
    "log"
)

// 解析配置文件

var (
    Setting *Config
)

const (
    DefaultBindAddress = "0.0.0.0:9277"
    DefaultBucket = 5;
    DefaultRedisHost = "127.0.0.1:6379"
    DefaultRedisDb = 1
    DefaultRedisPassword = ""
    DefaultRedisMaxIdle = 30
    DefaultRedisMaxActive = 0
)

type Config struct {
    BindAddress string
    Bucket int
    Redis RedisConfig
}

type RedisConfig struct {
    Host string
    Db int
    Password string
    MaxIdle int    // 连接池最大空闲连接数
    MaxActive int  // 连接池最大激活连接数
}

func Init(path string)  {
    Setting = &Config{}
    if (path == "") {
        Setting.initDefaultConfig()
        return
    }

    Setting.parse(path)
}

func (config *Config) parse(path string)  {
    file, err := ini.Load(path)
    if err != nil {
        log.Fatalf("无法解析配置文件#%s", err.Error())
    }

    section := file.Section("")
    config.BindAddress = section.Key("bind_address").MustString(DefaultBindAddress)
    config.Bucket = section.Key("bucket").MustInt(DefaultBucket)
    config.Redis.Host = section.Key("redis.host").MustString(DefaultRedisHost)
    config.Redis.Db = section.Key("redis.db").MustInt(DefaultRedisDb)
    config.Redis.Password = section.Key("redis.password").MustString(DefaultRedisPassword)
    config.Redis.MaxIdle = section.Key("redis.max_idle").MustInt(DefaultRedisMaxIdle)
    config.Redis.MaxActive = section.Key("redis.max_active").MustInt(DefaultRedisMaxActive)
}


func (config *Config) initDefaultConfig()  {
    config.BindAddress = DefaultBindAddress
    config.Bucket = DefaultBucket
    config.Redis.Host = DefaultRedisHost
    config.Redis.Db = DefaultRedisDb
    config.Redis.Password = DefaultRedisPassword
    config.Redis.MaxIdle = DefaultRedisMaxIdle
    config.Redis.MaxActive = DefaultRedisMaxActive
}