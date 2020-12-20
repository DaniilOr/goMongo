package cache

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

const cacheTimeout = 50 * time.Millisecond

type Service struct {
	Cache *redis.Pool
}

func (s *Service) FromCache(ctx context.Context, key string) ([]byte, error) {
	conn, err := s.Cache.GetContext(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	//id := ctx.Value(authenticator.AuthenticationContextKey)
	reply, err := redis.DoWithTimeout(conn, cacheTimeout, "GET", key)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	value, err := redis.Bytes(reply, err)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return value, err
}

func (s *Service) ToCache(ctx context.Context, key string, value []byte) error {
	conn, err := s.Cache.GetContext(ctx)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	//id := ctx.Value(authenticator.AuthenticationContextKey)
	_, err = redis.DoWithTimeout(conn, cacheTimeout, "SET", key, value)
	if err != nil {
		log.Print(err)
	}
	return err
}