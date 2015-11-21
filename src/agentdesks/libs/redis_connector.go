package libs

import (
	"agentdesks"
	"agentdesks/utils"
	"github.com/garyburd/redigo/redis"
)

func InitRedis(config *agentdesks.Config) redis.Conn {
	host := ":" + config.GetRedisPort()
	conn, err := redis.Dial("tcp", host)
	if err != nil {
		utils.PanicError(err, "Error establishing redis connection")
	}
	return conn
}