package models

import (
	"github.com/garyburd/redigo/redis"
)

type Connections struct {
	RedisConn redis.Conn
}
