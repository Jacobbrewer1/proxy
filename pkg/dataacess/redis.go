package dataacess

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	dbMonitoring "github.com/jacobbrewer1/reverse-proxy/pkg/dataacess/monitoring"
	"github.com/prometheus/client_golang/prometheus"
)

type IRedisDal interface {
	GetValue(key string) (string, error)
}

type redisDal struct {
	ctx    context.Context
	client *redis.Client

	collection int
}

func (r *redisDal) GetValue(key string) (string, error) {
	t := prometheus.NewTimer(dbMonitoring.RedisLatency.WithLabelValues(fmt.Sprintf("%d", r.collection)))
	defer t.ObserveDuration()

	data, err := r.client.WithContext(r.ctx).Get(key).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}

func NewRedisDal(ctx context.Context, collection int) IRedisDal {
	return &redisDal{
		ctx:        ctx,
		client:     Connections.RedisDb().Client(collection),
		collection: collection,
	}
}
