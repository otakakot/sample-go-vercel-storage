package main

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/redis/rueidis"
)

func main() {
	address := os.Getenv("KV_URL")

	opt, err := rueidis.ParseURL(address)
	if err != nil {
		panic(err)
	}

	opt.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	opt.DisableCache = true

	redis, err := rueidis.NewClient(opt)
	if err != nil {
		panic(err)
	}

	ping := redis.B().Ping().Build()

	ctx := context.Background()

	if err := redis.Do(ctx, ping).Error(); err != nil {
		panic(err)
	}

	set := redis.B().Set().Key("key").Value("test").ExSeconds(60).Build()

	if err := redis.Do(ctx, set).Error(); err != nil {
		panic(err)
	}

	get := redis.B().Get().Key("key").Build()

	val, err := redis.Do(ctx, get).ToString()
	if err != nil {
		panic(err)
	}

	println(val)

	del := redis.B().Del().Key("key").Build()

	if err := redis.Do(ctx, del).Error(); err != nil {
		panic(err)
	}
}
