package nosql

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
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
	address := config.Host + ":" + config.Port
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{address},
	})

	ctx := context.Background()
	pong, err := clusterClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect at", address)
		panic(err)
	}
	fmt.Println("Redis cluster", address, "connected!", pong)
	return &RedisCluster{
		clusterClient,
		ctx,
	}
}

type RedisCluster struct {
	clusterClient *redis.ClusterClient
	ctx           context.Context
}

func (redis *RedisCluster) Values(timeout int64, keys []string) ([]*Result, error) {
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
					"",
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

	go func() {
		wg.Wait()
		close(result)
	}()

	results := make([]*Result, len(keys))
	t := time.After(expiry)
	for {
		select {
		case r, open := <-result:
			if open {
				if r.err != nil {
					return nil, errors.New("redis Query Error")
				}
				results[r.Id] = r
			} else {
				fmt.Println("completed")
				return results, nil
			}
		case <-t:
			fmt.Println("timeOut")
			return nil, errors.New("time Out")
		}
	}
}

func (redis *RedisCluster) Shutdown() {

}
