package internal

import (
	"fmt"
	"simple_tiktok_rime/internal/dao/mysql"
	"simple_tiktok_rime/internal/dao/redis"
	_ "simple_tiktok_rime/internal/logic"
	"simple_tiktok_rime/internal/router"
	"simple_tiktok_rime/logs"
	"simple_tiktok_rime/manifest/config"
	"simple_tiktok_rime/pkg/cron"
	"simple_tiktok_rime/pkg/snowflake"

	"go.uber.org/zap"
)

const ConfigFilePath = "./manifest/config/config.yaml"

func Launch() {
	// 加载配置文件
	if err := config.Init(ConfigFilePath); err != nil {
		fmt.Printf("init config failed, Error:%v\n", err)
		return
	}

	// 加载日志配置
	if err := logs.Init(config.Conf.LogConfig); err != nil {
		fmt.Printf("init log failed, Error:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("log init success.")

	// 加载 mysql 配置
	if err := mysql.Init(config.Conf.MysqlConfig); err != nil {
		fmt.Printf("init mysql failed, Error:%v\n", err)
		return
	}
	defer mysql.Close()

	// 加载 redis 配置
	if err := redis.Init(config.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, Error:%v\n", err)
		return
	}
	defer redis.Close()

	// 加载雪花ID配置
	if err := snowflake.Init(config.Conf.StartTime, config.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, Error:%v\n", err)
		return
	}

	// 加载定时任务
	if err := cron.Init(); err != nil {
		fmt.Printf("init cron failed, Error:%v\n", err)
		return
	}

	// 启动服务
	router.Setup(router.Init())
}
