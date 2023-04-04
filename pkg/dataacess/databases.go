package dataacess

import "github.com/jacobbrewer1/reverse-proxy/pkg/dataacess/connection"

var Connections ConnectionSet

type ConnectionSet struct {
	redisDb *connection.RedisDb
}

func (c *ConnectionSet) RedisDb() *connection.RedisDb {
	return c.redisDb
}

func (c *ConnectionSet) SetRedisDb(redisDb *connection.RedisDb) {
	c.redisDb = redisDb
}
