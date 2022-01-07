package redis

import "time"

type Conf struct {
	RedisAddress string `xconf:"redis_address"`
}

type Timeout struct {
	ReadTimeout time.Duration
}
