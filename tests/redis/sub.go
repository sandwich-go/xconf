package redis

import "time"

// Conf test conf
type Conf struct {
	RedisAddress string `xconf:"redis_address"`
}

// Timeout test conf
type Timeout struct {
	ReadTimeout time.Duration
}
