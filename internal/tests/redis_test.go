package tests

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestRedisZRange(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			"127.0.0.1",
			6379,
		),
		DB:       0,
		PoolSize: 100,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		t.Fatal("connection redis failed", err)
	}

	videoFavoriteSet := "video_favorite"
	result, err := rdb.ZRange(context.Background(), videoFavoriteSet, 0, -1).Result()
	if err != nil {
		t.Fatal("err:", err)
	}
	fmt.Printf("result => %s", result)
}
