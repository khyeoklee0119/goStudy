package nosql

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type Config struct {
	Host       string
	Port       string
	FailFast   bool
	RetryDelay int
}

type Result struct {
	Id    int
	Key   string
	Value string
	err   error
}

type NoSQL interface {
	Values(...string) []string
	Shutdown() error
}

func NewRedisCluster(config Config) *RedisCluster {

	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{config.Host + config.Port},
	})

	ctx := context.Background()
	pong, err := clusterClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(pong)
	return &RedisCluster{
		clusterClient,
		ctx,
	}
}

type RedisCluster struct {
	clusterClient *redis.ClusterClient
	ctx           context.Context
}

func (redis *RedisCluster) Values(timeout int64, keys ...string) ([]*Result, error) {
	var wg sync.WaitGroup
	size := len(keys)
	wg.Add(size)
	result := make(chan *Result)
	expiry := time.Until(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	for id, key := range keys {
		go func(id int, key string) {
			defer wg.Done()
			value, err := redis.clusterClient.Get(context.Background(), key).Result()
			if err != nil {
				result <- &Result{
					id,
					key,
					nil,
					err,
				}
			}
			result <- &Result{
				id,
				key,
				value,
				nil,
			}
		}(id, key)
	}
	results := make([]*Result, len(keys))
	finish:=time.After(expiry)
	for {
		select {
		case r := <-result:
			results[r.Id] = r
		case <-finish:
			return nil, errors.New("Time Out")
		}
	}
	return results, nil
}

func (redis *RedisCluster) query(id int, key string, result chan *Result) {

}

func (redis *RedisCluster) Shutdown() {

}
