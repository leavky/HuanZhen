package util

import (
	"errors"
	"strconv"
	"strings"
)

// 解析IP格式, 80,81,8001-8004, 返回 ip 列表
func ParsePort(ips string) ([]int, error) {
	var portList []int
	ipStrList := strings.Split(ips, ",")
	for _, ipStrItem := range ipStrList {
		if strings.ContainsAny(ipStrItem, "-") {
			if len(strings.Split(ipStrItem, "-")) != 2 {
				return portList, errors.New("IP 格式错误")
			}
			startStrPort := strings.Split(ipStrItem, "-")[0]
			startPort, err := strconv.Atoi(startStrPort)
			if err != nil {
				return portList, errors.New("IP 格式错误")
			}
			endStrPort := strings.Split(ipStrItem, "-")[1]
			endPort, err := strconv.Atoi(endStrPort)
			if err != nil {
				return portList, errors.New("IP 格式错误")
			}
			if startPort >= endPort {
				return portList, errors.New("IP 格式错误")
			}

			for i := startPort; i <= endPort; i++ {
				portList = append(portList, i)
			}

		} else {
			port, err := strconv.Atoi(ipStrItem)
			if err != nil {
				return portList, err
			}
			portList = append(portList, port)
		}
	}

	// 去重
	portList = RemoveReplicaSliceInt(portList)
	return portList, nil
}

func RemoveReplicaSliceString(slc []string) []string {
	/*
		slice(string类型)元素去重
	*/
	result := make([]string, 0)
	tempMap := make(map[string]bool, len(slc))
	for _, e := range slc {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

func RemoveReplicaSliceInt(slc []int) []int {
	/*
		slice(int类型)元素去重
	*/
	result := make([]int, 0)
	tempMap := make(map[int]bool, len(slc))
	for _, e := range slc {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}
