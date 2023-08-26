# 项目技术栈

> gin + gorm + mysql + redis

# 项目结构说明

- api：接收前端传递的数据
- internal：内部应用具体逻辑，包含 router、controller、service、model 等业务模型信息
  - consts：包含常量信息
  - controller：接收 api 层数据并传递给业务层
  - dao：真正执行底层数据库操作的部分
  - logic：业务函数的实现部分
  - middleware：类似于 Java 中的拦截器概念
  - model：数据库与 go 结构体的一一对应，包含定义 dao 层的输入输出模型
  - router：应用的路由
  - service：业务函数的接口部分，具体的实现代码由 logic 来实现
  - tests：测试模块
  - launch.go： internal 内部的入口函数，只能由 main.go 来调用，该函数包含所有配置信息的初始化
- logs：记录应用产生的日志记录，以天为单位记录
- manifest：包含应用配置信息、部署配置信息、sql 文件等内容
- pkg 包含所有的外部工具类，包括但不限于第三方库、自己封装的工具等等

# 使用到的第三方库工具

- gin 框架：`go get -u "github.com/gin-gonic/gin"`
- gorm 框架：`go get -u "gorm.io/gorm"`
- go-redis ：`go get -u "github.com/go-redis/redis/v8"`
- viper 配置管理：`go get -u "github.com/spf13/viper"`
- air 热部署： `go install "github.com/cosmtrek/air"`
- zap 日志管理：
  - `go get -u "go.uber.org/zap"`
  - `go get -u "gopkg.in/natefinch/lumberjack.v2"`
- swagger 接口文档：`go get -u "github.com/go-swagger/go-swagger"`
- 七牛云的 GO-SDK：`go get github.com/qiniu/go-sdk/v7`
- 跨域处理：`go get github.com/gin-contrib/cors`
- 定时任务：`go get github.com/robfig/cron/v3`

