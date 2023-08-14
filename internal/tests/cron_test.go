package tests

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCorn(t *testing.T) {
	// 创建一个新的 Cron 实例
	c := cron.New()

	// 添加定时任务，使用 Cron 表达式来指定任务执行时间
	// 这里的示例会每分钟打印一条消息
	c.AddFunc("* * * * *", func() {
		fmt.Println("This job runs every minute.")
	})

	// 启动定时任务
	c.Start()

	// 程序运行一段时间来观察定时任务的执行
	time.Sleep(5 * time.Minute)

	// 停止定时任务
	c.Stop()
}
