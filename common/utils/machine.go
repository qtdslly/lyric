package utils

import (
	"net"
	"regexp"
	"strings"
)

//获取MAC地址
func GetMacAddr() string {
	result, err := BashCommand("ifconfig", 5)
	if err != nil {
		return ""
	}
	parttern := `[0-9a-z]{2}(:[0-9a-z]{2}){5}`
	re, _ := regexp.Compile(parttern)
	macs := re.FindAllString(result, 10)

	Mac := ""
	for _, network := range macs {
		Mac = Mac + network + ","
	}
	return strings.Trim(Mac, ",")
}

//获取IP地址
func GetIps() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	var ips []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	var result string
	for _, ip := range ips {
		result = result + "," + ip
	}

	return strings.Trim(string(result), ","), nil
}
