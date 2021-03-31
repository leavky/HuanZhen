package processCheck

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"log"
)

func StartProcessCheck() {
	log.Println("进程检测")

	fmt.Println("网络链接状态统计:")

	nc, _ := net.Connections("all")
	fmt.Println(nc)

	process, _ := process.Processes()
	for _, processItem := range process {
		processPid := processItem.Pid
		processName, _ := processItem.Name()
		processCwd, _ := processItem.Cwd()
		processConnections, _ := processItem.Connections()

		fmt.Println(processPid, processName, processCwd, processConnections)
	}

	//fmt.Println("进程统计:")
	//pi, _ := process.Pids()
	//fmt.Println(pi)
	//p, _ := process.NewProcess(614)
	//pm, _ := p.MemoryPercent()
	//pn, _ := p.Username()
	//fmt.Println(pm)
	//fmt.Println(pn)
}
