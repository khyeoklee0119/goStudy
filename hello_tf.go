package main

import "tensorflowGo/db/nosql"

func main() {
	nosql.NewRedisCluster(nosql.Config{Host: "dsp-dev-pctr.umxwaa.clustercfg.use1.cache.amazonaws.com", Port: "16379"})
}
