package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"monitor/config"
	"monitor/global"
	"monitor/initializers"
	"monitor/src/crontab"
	"os"
)

func main() {
	// 控制日志输出到文件
	f, _ := os.OpenFile("./log/app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	// 默认输出位置, 日志输出到文件和控制台两个位置
	gin.DefaultWriter = io.MultiWriter(f)
	router := gin.Default()
	config.InitConfig()
	global.DbMysqlXORM = initializers.InitMysqlCon()
	if global.DbMysqlXORM == nil {
		panic(any("数据库连接失败"))
	}
	defer global.DbMysqlXORM.Close()
	crontab.Cron()
	router.Run(":9001")
}
