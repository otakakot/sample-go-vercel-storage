package main

import (
	"context"
	"crypto/tls"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	address := os.Getenv("KV_URL")
	if address == "" {
		panic("KV_URL is required")
	}

	opt, err := redis.ParseURL(address)
	if err != nil {
		panic(err)
	}

	opt.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	cli := redis.NewClient(opt)

	if err := cli.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := cli.Set(ctx, "key", "test", time.Minute).Err(); err != nil {
		panic(err)
	}

	val, err := cli.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	println(val)

	if err := cli.Del(ctx, "key").Err(); err != nil {
		panic(err)
	}
}
