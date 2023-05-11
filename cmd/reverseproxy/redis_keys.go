package main

import (
	"errors"
	"fmt"
)

// redisKey is a custom type that allows us to define a key for a redis hash.
type redisKey string

const (
	// redisKeyRedirect is the key for the test redis hash where the redirect URL lives.
	redisKeyRedirect redisKey = "redirect"
)

// String implements fmt.Stringer
func (r redisKey) String() string {
	return string(r)
}

// RedisArg implements redis.Argument
func (r *redisKey) RedisArg() interface{} {
	if r == nil {
		return nil
	}
	return r.String()
}

// RedisScan implements redis.Scanner
func (r *redisKey) RedisScan(src interface{}) (err error) {
	if r == nil {
		return errors.New("nil pointer")
	}
	switch t := src.(type) {
	case []byte:
		if len(t) == 0 {
			return nil
		}
		*r = redisKey(t)
	default:
		err = fmt.Errorf("cannot convert from %T to %T", src, r)
	}
	return err
}
