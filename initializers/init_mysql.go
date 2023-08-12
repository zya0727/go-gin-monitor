package initializers

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"monitor/global"
)

func InitMysqlCon() *xorm.Engine {
	mysqlConfig := global.Config.Mysql
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Dbname)
	db, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		global.Log(global.LOG_ERROR, "initMysqlCon fail："+err.Error())
		return nil
	}
	db.ShowSQL(mysqlConfig.ShowSQL)
	return db
}
