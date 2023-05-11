package dataacess

import (
	"context"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	dbMonitoring "github.com/jacobbrewer1/reverse-proxy/pkg/dataacess/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type IRedisDal interface {
	GetValue(key string) (string, error)
}

type redisDal struct {
	ctx  context.Context
	conn redis.Conn

	collection int
}

func (r *redisDal) GetValue(key string) (string, error) {
	if dbMonitoring.RedisLatency != nil {
		t := prometheus.NewTimer(dbMonitoring.RedisLatency.WithLabelValues(fmt.Sprintf("%d", r.collection)))
		defer t.ObserveDuration()
	}
	got, err := r.conn.Do(redisGet, key)
	if err != nil {
		return "", err
	}
	if got == nil {
		return "", errors.New("key not found")
	}
	return string(got.([]byte)), nil
}

func NewRedisDal(collection int) IRedisDal {
	return NewRedisDalWithContext(context.Background(), collection)
}

func NewRedisDalWithContext(ctx context.Context, collection int) IRedisDal {
	if ctx == nil {
		ctx = context.Background()
	}

	conn, err := Connections.RedisDb().Conn(collection)
	if err != nil {
		log.Println("Error getting redis connection", err)
		return nil
	}

	return &redisDal{
		ctx:        ctx,
		conn:       conn,
		collection: collection,
	}
}
