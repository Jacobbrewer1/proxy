package connection

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisDb struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     string `json:"port,omitempty" yaml:"port,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// clients is a map of the database number to the client for that database
	clients map[int]*redis.Client
}

func (r *RedisDb) Client(database int) *redis.Client {
	if _, ok := r.clients[database]; !ok {
		r.generateClient(database)
	}

	client, _ := r.clients[database]

	client.Ping()
	return client
}

func (r *RedisDb) generateClient(database int) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
		DB:       database,
	})

	if r.clients == nil {
		r.clients = make(map[int]*redis.Client, 16)
	}

	r.clients[database] = client
}
