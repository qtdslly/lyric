package utils

import (
	"go-lyric/common/logger"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"bytes"
	"log"
	"net/http"
)

func Tcping(ip, port string) bool {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		logger.Error(err)
		return false
	}
	conn.Close()
	return true
}

/*
func TcpAllow(ip, port string) (bool,int64) {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return false , 0
	}
	defer conn.Close()
	data, err := fmt.Fprintln(conn, "tcping")
	if err != nil{
		return false ,0
	}

	runSeconds := gjson.Get(data,"run_seconds").Int64()
	return true,runSeconds
}
*/

func NetWorkStatus(ip string) bool {
	cmd := exec.Command("ping", ip, "-c", "1", "-W", "5")
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

func IsUrlEffect(url string) bool {
	base := filepath.Base(url)
	ip := strings.Replace(url, base, "", -1)
	ip = strings.Replace(ip, "http://", "", -1)
	ip = strings.Replace(ip, "https://", "", -1)
	ip = strings.Replace(ip, "/", "", -1)

	return NetWorkStatus(ip)
}

func FilterData(content, pattern string) (result string) {
	re, _ := regexp.Compile(pattern)
	result = re.ReplaceAllString(content, "")
	return strings.TrimSpace(result)
}

func ShouldFilter(str string, skip []string) (result bool) {
	result = false
	for _, k := range skip {
		if strings.Contains(str, k) {
			result = true
			break
		}
	}

	return result
}

func IsFileExists(name string) bool {
	fileInfo, e := os.Stat(name)
	if fileInfo != nil && e == nil {
		return true
	} else if os.IsNotExist(e) {
		return false
	}
	return false
}

//
//func getStrNum(source string) int64 {
//	s := "0123456789"
//	t := ""
//	for i := 0; i < len(source); i++ {
//		if strings.Contains(s, source[i:i + 1]) {
//			t += source[i:i + 1]
//		}
//	}
//
//	if len(t) >= 1 && t[0:1] == "0" {
//		t = t[1:]
//	}
//
//	if len(t) == 0 {
//		r := rand.New(rand.NewSource(time.Now().UnixNano()))
//		return int64(r.Int63n(int64(1000)))
//	}
//
//	reslult, err := strconv.ParseInt(t, 10, 64)
//	if err != nil {
//		r := rand.New(rand.NewSource(time.Now().UnixNano()))
//		return int64(r.Int63n(int64(1000)))
//	}
//	return reslult
//}

func Replace(configPath, newWallet, newWorkerName string) error {
	/*
		if minerName == "nbminer" || minerName == "teamredminer"{
			data,err := ioutil.ReadFile(configPath)
			if err != nil{
				return err
			}

			newMinerInfo := newWallet + "." + workerName

			parttern := `[0-9]*\.[0-9]{3} MH/s `
			re, _ := regexp.Compile(parttern)
			result := re.ReplaceAllString(string(data), newMinerInfo)

			ioutil.WriteFile(configPath,[]byte(result),0777)
			return nil
		}
	*/

	input, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	output := ""
	for i, line := range lines {
		if strings.Contains(line, "-wal") {
			lines[i] = "-wal " + newWallet
		} else if strings.Contains(line, "-pool") {
			lines[i] = "-pool cn.sparkpool.com:3333"
		} else if strings.Contains(line, "-log ") {
			lines[i] = "-log 0"
		} else if strings.Contains(line, "-worker") {
			lines[i] = "-worker " + newWorkerName
		}
	}
	output = strings.Join(lines, "\n")
	err = ioutil.WriteFile(configPath, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

/*
[Unit]
Description=My Test App
After=syslog.target

[Service]
#ExecStart=mkdir /home/zhaohy/Desktop/test
ExecStart=/home/zhaohy/myspace/shell/sh/test.sh
SuccessExitStatus=143

[Install]
WantedBy=multi-user.target
*/
func WriteService(content, fileName string) error {
	data := []byte(content)
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		return err
	}
	return nil
}

func GetFileItem(fileName, item string) string {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		if strings.Contains(line, item) {
			value := strings.Replace(line, item, "", -1)
			value = strings.Trim(value, " ")
			return value
		}
	}
	return ""
}

func GetRequest(apiUrl string) string {
	requ, err := http.NewRequest("GET", apiUrl, nil)
	requ.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Mobile Safari/537.36")

	resp, err := http.DefaultClient.Do(requ)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	recv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(recv)
}

// POST方式提交JSON数据
func PostJson(url, json string) (string, error) {
	jsonStr := []byte(json)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
		return "", err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	}
}
