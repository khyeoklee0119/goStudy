package main

import (
	"fmt"
	"tensorflowGo/db/nosql"
)

func main() {
	redis := nosql.NewRedisCluster(nosql.Config{Host: "127.0.0.1", Port: "30001"})
	keys := []string{"test1", "test2", "test3", "test4", "test5", "test6"}
	results, err := redis.Values(5, keys)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Println("---------Result ---------")
	for _, result := range results {
		fmt.Println(result.Id+1, result.Key, ":", result.Value)
	}
}
