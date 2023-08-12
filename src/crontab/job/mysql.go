package job

import (
	"errors"
	"monitor/global"
	"strconv"
	"strings"
)

var mysqlMonitor global.MysqlMonitor

func MysqlMonitor() {
	mysqlMonitor = global.Config.MysqlMonitor

	var errorMsg []string
	_, conErr := connectionsNum()
	if conErr != nil {
		errorMsg = append(errorMsg, conErr.Error())
	}
	_, transactionError := longTransaction()
	if transactionError != nil {
		errorMsg = append(errorMsg, transactionError.Error())
	}
	_, dirtyPageErr := dirtyPageRate()
	if dirtyPageErr != nil {
		errorMsg = append(errorMsg, dirtyPageErr.Error())
	}
	threadNum()

	if len(errorMsg) > 0 {
		//TODO 邮件
		global.Log(global.LOG_ERROR, "[MYSQL]"+strings.Join(errorMsg, "\n"))
	}
}

/**
并发连接数
*/
func connectionsNum() (bool, error) {
	res, err := global.DbMysqlXORM.QueryString("show status like 'Threads_connected'")
	if err != nil {
		return false, err
	}
	currentDbConNum, _ := strconv.Atoi(res[0]["Value"])
	maxCon, err2 := global.DbMysqlXORM.QueryString("show variables like '%max_connections%'")
	if err2 != nil {
		return false, err2
	}
	var maxAllowNum int
	for _, varInfo := range maxCon {
		if strings.EqualFold(varInfo["Variable_name"], "max_connections") {
			maxAllowNum, _ = strconv.Atoi(varInfo["Value"])
		}
	}
	//max_connections:当前数据库能够连接的所有连接数，不区分用户。
	//max_user_connections:当前数据库用户所能连接的连接数，不会超过max_connections。
	//超过设置的比例 比如90% 警告
	rate := mysqlMonitor.MaxConnectionsRate
	monitorMaxCon := int(float64(maxAllowNum) * (float64(rate) / 100.0))
	if currentDbConNum >= monitorMaxCon {
		return false, errors.New("连接数超过允许的连接数的" + strconv.Itoa(rate) + " %（" + strconv.Itoa(monitorMaxCon) + "） 当前连接数：" + strconv.Itoa(currentDbConNum))
	}
	return true, nil
}

/**
长事务
*/
func longTransaction() (bool, error) {
	transactionMaxAllowSecond := mysqlMonitor.TransactionMxAllowSecond
	res, err := global.DbMysqlXORM.QueryString("select count(*) as cnt from information_schema.innodb_trx where TIME_TO_SEC(timediff(now(),trx_started))>" + strconv.Itoa(transactionMaxAllowSecond))
	if err != nil {
		return false, errors.New("查询事务时异常：" + err.Error())
	}
	longTransactionNum, _ := strconv.Atoi(res[0]["cnt"])
	if longTransactionNum >= 1 {
		return false, errors.New("目前有超过" + strconv.Itoa(transactionMaxAllowSecond) + "s的" + strconv.Itoa(longTransactionNum) + "条长事务查询在执行。")
	}
	return true, nil
}

func dirtyPageRate() (bool, error) {
	//不要让脏页比例超过 75%
	/**
	  5.7 以后 performance_schema.global_status
	  5.7之前 information_schema.global_status(8.0被废弃）
	由于在代码中使用了github.com/go-sql-driver/mysql ，一直出现语法错误，后来定位到是在一个sql语句中执行 multi statements
	默认是不支持multi statements的需要进行配置，因为 multi statements 可能会增加sql注入的风险
	*/
	dirtyPage, err := global.DbMysqlXORM.QueryString("select VARIABLE_VALUE from `performance_schema`.global_status where VARIABLE_NAME = 'Innodb_buffer_pool_pages_dirty'")
	if err != nil {
		return false, err
	}
	dirtyPageInt, _ := strconv.Atoi(dirtyPage[0]["VARIABLE_VALUE"])
	dirtyPageTotal, totalErr := global.DbMysqlXORM.QueryString("select VARIABLE_VALUE from `performance_schema`.global_status where VARIABLE_NAME = 'Innodb_buffer_pool_pages_total'")
	if totalErr != nil {
		return false, totalErr
	}
	dirtyPageTotalInt, _ := strconv.Atoi(dirtyPageTotal[0]["VARIABLE_VALUE"])

	rate := dirtyPageInt * 100 / dirtyPageTotalInt
	if rate >= mysqlMonitor.DirtyPageMaxRate {
		return false, errors.New("数据库脏页比例达到" + strconv.Itoa(mysqlMonitor.DirtyPageMaxRate) + "%以上，会影响查询效率。")
	}
	return true, nil
}

func threadNum() {
	// innodb_thread_concurrency 控制并发线程数量 默认值是0 表示不限制并发线程数量
}
