package job

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"monitor/global"
	"strings"
)

var serverMonitor global.ServerMonitor

func ServerMonitor() {
	var errorMsg []string
	serverMonitor = global.Config.ServerMonitor
	checkCpu()
	memErr := checkMen()
	if memErr != nil {
		errorMsg = append(errorMsg, memErr.Error())
	}
	if len(errorMsg) > 0 {
		//TODO 邮件
		global.Log(global.LOG_ERROR, "[SERVER]"+strings.Join(errorMsg, "\n"))
	}
}

func checkMen() error {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return errors.New("获取内存失败：" + err.Error())
	}
	memMonitorPercent := serverMonitor.MemMaxUsedPercent
	if memInfo.UsedPercent >= memMonitorPercent {
		return errors.New(fmt.Sprintf("内存使用比例过大：%f %", memInfo.UsedPercent))
	}
	return nil
}

func checkCpu() {
	/**top 只能持续运行一段时间，而 ps 是立刻返回的。
	这个差异体现在运行top -n 1和ps aux时，top是延迟后返回的，而ps是立刻返回的。
	*/
}
