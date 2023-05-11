package connection

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"net"
	"time"
)

type RedisDb struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     string `json:"port,omitempty" yaml:"port,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// conns is a map of the database number to the client for that database
	conns map[int]redis.Conn
}

func (r *RedisDb) Conn(database int) (redis.Conn, error) {
	conn, ok := r.conns[database]
	if !ok {
		if err := r.generateClient(database); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (r *RedisDb) generateClient(database int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addr := net.JoinHostPort(r.Host, r.Port)
	client, err := redis.DialContext(ctx, "tcp", addr, redis.DialPassword(r.Password), redis.DialDatabase(database))
	if err != nil {
		return err
	}

	if r.conns == nil {
		r.conns = make(map[int]redis.Conn, 16)
	}
	r.conns[database] = client
	return nil
}
