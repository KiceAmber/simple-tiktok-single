package redis

import (
	"context"
	"fmt"
	"simple_tiktok_rime/manifest/config"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client
	once = sync.Once{}
)

func New() *redis.Client {
	return rdb
}

func Init(cfg *config.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		//Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping(context.Background()).Result()
	return
}

func Close() {
	_ = rdb.Close()
}

// ================ redis 相关的操作 ================

// SetIsMember 判断 member 是否在 Set 中
//func SetIsMember(ctx context.Context, key string, member any) (bool, error) {
//	return rdb.SIsMember(ctx, key, member).Result()
//}

// ZSetIncrNumber 给 ZSet 的 key 对应的 score 增加 number
//func ZSetIncrNumber(ctx context.Context, key string, number float64, member string) error {
//	return rdb.ZIncrBy(ctx, key, number, member).Err()
//}

// SetAddMember Set 增加 member
//func SetAddMember(ctx context.Context, key string, member any) error {
//	return rdb.SAdd(ctx, key, member).Err()
//}

// 移除 key 对应的 member
//func SetRemoveMember(ctx context.Context, key string, member any) error {
//	return rdb.SRem(ctx, key, member).Err()
//}
