package global

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"time"
)

type AppConfig struct {
	App           App           `json:"app"`
	Mysql         Mysql         `json:"mysql"`
	MysqlMonitor  MysqlMonitor  `mapstructure:"mysql-monitor" json:"mysql-monitor" yaml:"mysql-monitor"` //mapstructure 这个要加不然读不出来
	ServerMonitor ServerMonitor `mapstructure:"server-monitor" json:"mysql-monitor" yaml:"mysql-monitor"`
}

type App struct {
	Port int `json:"port"`
}

type Mysql struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Dbname   string `json:"dbname,omitempty"`
	ShowSQL  bool   `mapstructure:"show-sql" json:"show-sql,omitempty" yaml:"show-sql"`
	IsLocal  bool   `mapstructure:"is-local" json:"is-local,omitempty" yaml:"is-local"`
}

type MysqlMonitor struct {
	MaxConnectionsRate       int `mapstructure:"max-connections-rate" json:"max-connections-rate,omitempty" yaml:"max-connections-rate"`
	TransactionMxAllowSecond int `mapstructure:"transaction-max-allow-second" json:"transaction-max-allow-second,omitempty" yaml:"transaction-mx-allow-second"`
	DirtyPageMaxRate         int `mapstructure:"dirty-page-max-rate" json:"dirty-page-max-rate,omitempty" yaml:"dirty-page-max-rate"`
	MaxCpu                   int `mapstructure:"max-cpu" json:"max-cpu,omitempty" yaml:"max-cpu"`
}

type ServerMonitor struct {
	MemMaxUsedPercent float64 `mapstructure:"mem-max-used-percent" json:"mem-max-used-percent,omitempty" yaml:"mem-max-used-percent"`
}

var (
	DbMysqlXORM *xorm.Engine
	Config      AppConfig
)

const (
	LOG_INFO  = "INFO"
	LOG_ERROR = "ERROR"
)

func Log(logType string, logMsg string) {
	now := time.Now().Format("2006/01/02 - 15:04:05")
	f := fmt.Sprintf("[%s] %s %s\n", logType, now, logMsg)
	fmt.Fprintf(gin.DefaultWriter, f)
}
