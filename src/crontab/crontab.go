package crontab

import (
	"github.com/robfig/cron"
	"monitor/global"
	"monitor/initializers"
	"monitor/src/crontab/job"
)

func Cron() {
	cron := cron.New()
	cron.AddFunc("1 * * * *", func() {
		err := global.DbMysqlXORM.Ping()
		if err != nil {
			global.Log(global.LOG_ERROR, "database ping fail:"+err.Error())
			//重连
			global.DbMysqlXORM = initializers.InitMysqlCon()
			rePingErr := global.DbMysqlXORM.Ping()
			if rePingErr != nil {
				global.Log(global.LOG_ERROR, " reping fail"+rePingErr.Error())
				return
			}
		}
		global.Log(global.LOG_INFO, "mysql monitor start")
		job.MysqlMonitor()
		global.Log(global.LOG_INFO, "mysql monitor end")
		global.Log(global.LOG_INFO, "server monitor start")
		job.ServerMonitor()
		global.Log(global.LOG_INFO, "server monitor end")
	})
	cron.Start()
}
