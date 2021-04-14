package processCheck

import (
	"HuanZhen/logger"
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var blackIpList []string
var blackFileMd5 []string

func StartProcessCheck() {
	log.Println("启动进程检测")
	// load black ip
	fi, err := os.Open("./Asset/BlackIp/blackIp.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		blackIpList = append(blackIpList, string(a))
	}

	// load black file md5
	fi2, err := os.Open("./Asset/Malware-Md5/Malware-Md5.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi2.Close()

	br2 := bufio.NewReader(fi)
	for {
		a, _, c := br2.ReadLine()
		if c == io.EOF {
			break
		}
		md5Str := strings.Split(string(a), " ")[0]
		blackFileMd5 = append(blackFileMd5, md5Str)
	}

	checkProcess()
}

func checkProcess(){
	for {
		process, _ := process.Processes()
		for _, processItem := range process {
			processPid := processItem.Pid
			processConnections, _ := processItem.Connections()
			processExe,_ := processItem.Exe()

			// check process connection
			for _, connItem := range processConnections{
				remoteIp := connItem.Raddr.IP
				if MatchBlackIp(remoteIp){
					logger.HZLogger.Warn("检测到有进程连接恶意IP", processPid, "--", processExe)
				}
			}
			// check process exe md5
			if processExe != ""{
				log.Println(processExe)
				fileMd5 := getFileMd5(processExe)
				if MatchBlackMd5(fileMd5) {
					logger.HZLogger.Info("检测到系统中存在恶意文件", processPid, "--", processExe)
				}
			}
		}

		// 1 minute
		time.Sleep(time.Minute * 1)
	}

}

func MatchBlackIp(ip string)bool{
	for _, ipItem := range blackIpList{
		if ipItem == ip{
			return true
		}
	}
	return false
}

func MatchBlackMd5(fileMd5 string)bool{
	for _, fileMd5Item := range blackFileMd5{
		if fileMd5Item == fileMd5{
			return true
		}
	}
	return false
}

func getFileMd5(filePath string)string{
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
