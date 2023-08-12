package crontab

import (
	"github.com/robfig/cron"
	"monitor/global"
	"monitor/src/crontab/job"
)

func Cron() {
	cron := cron.New()
	cron.AddFunc("1 * * * *", func() {
		global.Log(global.LOG_INFO, "mysql monitor start")
		job.MysqlMonitor()
		global.Log(global.LOG_INFO, "mysql monitor end")
		global.Log(global.LOG_INFO, "server monitor start")
		job.ServerMonitor()
		global.Log(global.LOG_INFO, "server monitor end")
	})
	cron.Start()
}
