package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/tidwall/gjson"
)

func GetAmdGPUCount() int {
	var args []string
	result, err := Command(Decrypt("f676034afea3fd398a91bddfcea1bd7f84e6849f3bfd920ad8a985c08aafb4460184618c45182ebddb6f"), 5, args)
	if err != nil {
		return 0
	}

	n, err := strconv.Atoi(strings.Trim(result, "\n"))
	if err != nil {
		return 0
	}
	return n
}

type MiningResult struct {
	data string
	err  error
}

// 通过接口取挖矿统计
func getMiningStatic(miner Miner, ch chan bool, result *MiningResult) {
	mi, err := miner.GetInfo()
	if err != nil {
		result.data, result.err = "", err
	}

	result.data, result.err = mi.Json()
	ch <- true
}
func GetMiningInfo() (string, error) {
	miner := Miner{Address: "127.0.0.1:3333", Password: ""}
	t := time.NewTimer(time.Duration(3) * time.Second)
	done := make(chan bool)
	result := &MiningResult{data: "", err: nil}
	go getMiningStatic(miner, done, result)

	select {
	case <-done:
		return string(result.data), result.err
	case <-t.C:
		result.err = errors.New("timeout")
		return string(result.data), result.err
	}
}
func FloatToStr(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

func GetCurrentHashrate() float64 {
	cmd := Decrypt("78ea6edb4f0ad61ecb9b393902f76741a745309735f8242148b1e1eb65639ec5933af9ebdb29db333f79a8d92c8690f6592aaaa90c1cbe88ce60c5cc364e67e3526585918995f6a4462e9bc65642dddc43a0996c2f4205c862bfc4cf0b19b0838d79c1273f4105712e883538e8134aef12")
	data, err := BashCommand(cmd, 5)
	if err != nil {
		return 0.0
	}
	data = strings.Replace(data, "\n", "", -1)
	hashrate, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return 0.0
	}
	return hashrate
}

func GetCurrentHashrate1(kernel string) float64 {
	if kernel == "nbminer" {
		return GetNbminer()
	} else if kernel == "phoenixminer" {
		return GetPhoenixMiner()
	} else if kernel == "trexminer" {
		return GetTrexMinerMiningData()
	} else if kernel == "lolminer" {
		return GetlolMinerMiningData()
	}

	return 0.0
}
func GetPhoenixMiner() float64 {
	MinInfo, err := GetMiningInfo()
	if err != nil {
		return 0.0
	}
	MinJson, err := gabs.ParseJSON([]byte(MinInfo))
	if err != nil {
		return 0.0
	}

	hashrates := FloatToStr(MinJson.Path("MainCrypto.HashRate").Data().(float64) / 1000)
	hashrate, err := strconv.ParseFloat(hashrates, 64)
	if err != nil {
		return 0.0
	}

	return hashrate
}

func GetNbminer() float64 {
	hashrate := 0.0
	url := "http://127.0.0.1:3334/api/v1/status"
	result, err := GetReq(url)
	if err != nil {
		return hashrate
	}
	miningInfo := gjson.Parse(result)
	sHashrate := strings.Trim(strings.Replace(miningInfo.Get("miner.total_hashrate").String(), "M", "", -1), " ")
	hashrate, _ = strconv.ParseFloat(sHashrate, 64)
	return hashrate
}

func GetlolMinerMiningData() float64 {
	hashrate := 0.0
	url := "http://127.0.0.1:3336"
	result, err := GetReq(url)
	if err != nil {
		return hashrate
	}

	miningInfo := gjson.Parse(result)

	sHashrate := miningInfo.Get("Session.Performance_Summary").String()
	hashrate, _ = strconv.ParseFloat(sHashrate, 64)

	return hashrate
}

func GetTrexMinerMiningData() float64 {
	hashrate := 0.0
	url := "http://127.0.0.1:3337/summary"
	result, err := GetReq(url)
	if err != nil {
		return hashrate
	}

	miningInfo := gjson.Parse(result)
	hashrate = float64(miningInfo.Get("hashrate").Int()) / float64(1000000)

	return hashrate
}

func GetReq(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		return string(body), nil
	}
}

var args = struct {
	id      string // 命令ID
	jsonrpc string
	psw     string
}{"0", "2.0", ""}

// -------------------------------------------------------------------------------------------
// 算力相关
type Crypto struct {
	HashRate       int
	Shares         int
	RejectedShares int
	InvalidShares  int
}

// -------------------------------------------------------------------------------------------
func (c Crypto) String() (s string) {
	if c.HashRate+c.RejectedShares+c.InvalidShares+c.Shares == 0 {
		return "Disabled\n"
	}
	s += fmt.Sprintf("HashRate:         %8d Mh/s\n", c.HashRate)
	s += fmt.Sprintf("Shares:           %8d\n", c.Shares)
	s += fmt.Sprintf("Rejected Shares:  %8d\n", c.RejectedShares)
	s += fmt.Sprintf("Invalid Shares:   %8d\n", c.InvalidShares)
	return s
}

// 自定义错误消息
type opError struct {
	s string
}

// -------------------------------------------------------------------------------------------
// 算力的JSON格式
func (c Crypto) Json() (string, error) {
	if c.HashRate+c.RejectedShares+c.InvalidShares+c.Shares == 0 {
		return "", errors.New("Disabled")
	}

	result, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// -------------------------------------------------------------------------------------------
// 矿池信息
type PoolInfo struct {
	Address  string
	Switches int
}

// -------------------------------------------------------------------------------------------
// 矿池信息
func (p PoolInfo) String() (s string) {
	if p.Address == "" {
		return "Disabled\n"
	}

	s += fmt.Sprintf("Address:   %s\n", p.Address)
	s += fmt.Sprintf("Switches:  %"+strconv.Itoa(len(p.Address))+"d\n", p.Switches)
	return s
}

// -------------------------------------------------------------------------------------------
func (p PoolInfo) Json() (string, error) {
	if p.Address == "" {
		return "", errors.New("Disabled")
	}

	result, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// -------------------------------------------------------------------------------------------
// 单个GPU 信息
type GPU struct {
	HashRate    int
	AltHashRate int
	Temperature int
	FanSpeed    int
}

// -------------------------------------------------------------------------------------------
// GPU 信息转为字符串
func (gpu GPU) String() (s string) {
	s += fmt.Sprintf("Hash Rate:     %8d Mh/s\n", gpu.HashRate)
	s += fmt.Sprintf("Alt Hash Rate: %8d Mh/s\n", gpu.AltHashRate)
	s += fmt.Sprintf("Temperature:   %8d º\n", gpu.Temperature)
	s += fmt.Sprintf("Fan Speed:     %8d %%\n", gpu.FanSpeed)
	return s
}

func (gpu GPU) Json() (string, error) {
	result, err := json.Marshal(gpu)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// -------------------------------------------------------------------------------------------
// 检测GPU是否挂掉
// return [true:挂掉，没有挖矿， false: 正常工作，可以挖矿]
func (gpu GPU) IsStuck() bool {
	return gpu.HashRate == 0
}

// -------------------------------------------------------------------------------------------
// 挖矿信息
type MinerInfo struct {
	Version    string   // 版本
	UpTime     int      // 运行时间(单位：分)
	MainCrypto Crypto   // 主挖币种
	AltCrypto  Crypto   // 双挖币种
	MainPool   PoolInfo // 首选矿池
	AltPool    PoolInfo // 备用矿池
	GPUS       []GPU    // GPU 列表
}

// -------------------------------------------------------------------------------------------
// 统计已经卡死的GPU数量, 即没有挂矿的GPU数量
func (m MinerInfo) StuckGPUs() int {
	var total int
	for _, gpu := range m.GPUS {
		if gpu.IsStuck() {
			total++
		}
	}
	return total
}

// -------------------------------------------------------------------------------------------
// 字符串形式的矿工信息
func (m MinerInfo) String() string {
	var s string
	s += fmt.Sprintf("Version:   %10s\n", m.Version)
	s += fmt.Sprintf("Up Time:   %10d min\n", m.UpTime)
	s += "\n"
	s += fmt.Sprintf("Main Crypto\n%s\n", m.MainCrypto)
	s += fmt.Sprintf("Alt Crypto\n%s\n", m.AltCrypto)
	s += fmt.Sprintf("Main Pool\n%s\n", m.MainPool)
	s += fmt.Sprintf("Alt Pool\n%s\n", m.AltPool)
	for i, gpu := range m.GPUS {
		s += fmt.Sprintf("GPU %d\n%s\n", i, gpu)
	}
	return s
}

func (m MinerInfo) Json() (string, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// -------------------------------------------------------------------------------------------
// 矿工
type Miner struct {
	Address  string
	Password string
}

const (
	methodGetInfo      = "miner_getstat1" // 获取统计信息
	methodRestartMiner = "miner_restart"  // 重启挖矿程序
	methodReboot       = "miner_reboot"   // 重启电脑
	methodControlGPU   = "control_gpu"    // GPU设置
	methodGetConfig    = "miner_getfile"  // 获取配置文件信息 config.txt
	methodSetConfig    = "miner_file"     // 更新配置文件	config.txt
)

// -------------------------------------------------------------------------------------------
// 矿工地址
func (m Miner) String() (s string) {
	return fmt.Sprintf("Miner {Address: %s}\n", m.Address)
}

// -------------------------------------------------------------------------------------------
// 重启挖矿程序
func (m Miner) Restart() error {
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	args.psw = m.Password
	return client.Call(methodRestartMiner, args, nil)
}

// -------------------------------------------------------------------------------------------
// 重启电脑
func (m Miner) Reboot() error {
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	args.psw = m.Password
	return client.Call(methodReboot, args, nil)
}

// -------------------------------------------------------------------------------------------
// 单个GPU状态控制
// index : GPU 序号 -1 全部GPU
// status: GPU的新状态, 0 - 禁用, 1 - 仅 ETH , 2 - 双挖模式
func (m Miner) ControlGPU(index, status int) error {
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	args.psw = m.Password
	var params = struct {
		id      string // 命令ID
		jsonrpc string
		psw     string
		params  []int
	}{args.id, args.jsonrpc, args.psw, []int{index, status}}

	return client.Call(methodControlGPU, params, nil)
}

// -------------------------------------------------------------------------------------------
// 获取配置文件 config.txt
func (m Miner) GetConfig() (string, error) {
	var reply string
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return "", err
	}
	defer client.Close()

	args.psw = m.Password

	var params = struct {
		id      string // 命令ID
		jsonrpc string
		psw     string
		params  []string
	}{args.id, args.jsonrpc, args.psw, []string{"config.txt"}}

	err = client.Call(methodGetConfig, params, &reply)
	if err != nil {
		return "", err
	}

	return reply, nil
}

// -------------------------------------------------------------------------------------------
// 更新配置文件
func (m Miner) UpdateConfig(content string) error {
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	args.psw = m.Password

	var params = struct {
		id      string // 命令ID
		jsonrpc string
		psw     string
		params  []string
	}{args.id, args.jsonrpc, args.psw, []string{"config.txt", content}}

	return client.Call(methodGetConfig, params, nil)
}

// -------------------------------------------------------------------------------------------
// 获取Claymore或PhoenixMiner 当前统计信息
func (m Miner) GetInfo() (MinerInfo, error) {
	var mi MinerInfo
	var reply []string
	client, err := jsonrpc.Dial("tcp", m.Address)
	if err != nil {
		return mi, err
	}
	defer client.Close()

	args.psw = m.Password
	err = client.Call(methodGetInfo, args, &reply)
	if err != nil {
		return mi, err
	}

	return parseResponse(reply), nil
}

// -------------------------------------------------------------------------------------------
// 解析统计结果
func parseResponse(info []string) MinerInfo {
	var mi MinerInfo
	var group []string

	mi.Version = strings.Replace(info[0], " - ETH", "", 1) //版本
	mi.UpTime = toInt(info[1])                             //在线时长

	//总算力
	group = splitGroup(info[2])
	mi.MainCrypto.HashRate = toInt(group[0])
	mi.MainCrypto.Shares = toInt(group[1])
	mi.MainCrypto.RejectedShares = toInt(group[2])

	// DCR
	group = splitGroup(info[4])
	mi.AltCrypto.HashRate = toInt(group[0])
	mi.AltCrypto.Shares = toInt(group[1])
	mi.AltCrypto.RejectedShares = toInt(group[2])

	//主挖矿池
	group = splitGroup(info[7])
	mi.MainPool.Address = group[0]
	if len(group) > 1 {
		mi.AltPool.Address = group[1]
	}

	// eth
	group = splitGroup(info[8])
	mi.MainCrypto.InvalidShares = toInt(group[0]) //umber of ETH invalid shares
	mi.MainPool.Switches = toInt(group[1])        //number of ETH pool switches
	mi.AltCrypto.InvalidShares = toInt(group[2])  //number of DCR invalid shares
	mi.AltPool.Switches = toInt(group[3])         //number of DCR pool switches.

	// 每GPU的算力
	for _, hashrate := range splitGroup(info[3]) {
		mi.GPUS = append(mi.GPUS, GPU{HashRate: toInt(hashrate)})
	}

	// GUP 温度，风扇转速
	for i, val := range splitGroup(info[6]) {
		if i%2 == 0 {
			mi.GPUS[i/2].Temperature = toInt(val)
		} else {
			mi.GPUS[(i-1)/2].FanSpeed = toInt(val)
		}
	}

	// 备用矿池
	if mi.AltPool.Address != "" {
		for i, val := range splitGroup(info[5]) {
			hashrate, err := strconv.Atoi(val)
			if err == nil {
				mi.GPUS[i].AltHashRate = hashrate
			}
		}
	}

	return mi
}

// -------------------------------------------------------------------------------------------
// 解析每行数据
func splitGroup(s string) []string {
	return strings.Split(s, ";")
}

// -------------------------------------------------------------------------------------------
// 字符串转数值
func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

type GpuStatics struct {
	Name        string //名称
	Power       string //功耗
	MemClock    string //显存频率
	BusId       string //pci bus_id
	MaxPower    string //GPU 功耗限制 适用于nvidia 卡
	Type        string //GPU 类型, A, N
	Driver      string //驱动版本 适用nvidia
	Memory      string //显存 适用nvidia
	VbVersion   string //bios版本
	Brand       string //显存品牌型号
	Subvendor   string //厂家
	CoreClock   string //核心频率
	Temperature string //温度
	FanSpeed    string //风扇
}

func GetGPUInfo() (map[int]GpuStatics, error) {
	m := make(map[int]GpuStatics)
	rep, err := Command("/home/ymminer/bin/gpu-detect", 30, GetArgs("listjson"))
	if err == nil {
		body := rep[strings.Index(rep, "[") : len(rep)-1]

		i := 0

		list := gjson.Parse(body)
		items := list.Array()
		for _, item := range items {
			var gpu GpuStatics
			var card string

			brand := item.Get("brand").String()
			if brand == "amd" {
				card = "A"
			} else if brand == "nvidia" {
				card = "N"
			} else {
				continue //跳过内置GPU（brand=cpu）
			}

			if item.Get("name").Exists() {
				gpu.Name = item.Get("name").String()
			}
			if item.Get("busid").Exists() {
				gpu.BusId = item.Get("busid").String()
			}
			gpu.Type = card

			if item.Get("subvendor").Exists() {
				gpu.Subvendor = item.Get("subvendor").String()
			}
			if item.Get("vbios").Exists() {
				gpu.VbVersion = item.Get("vbios").String()
			}
			if item.Get("mem_type").Exists() {
				gpu.Brand = item.Get("mem_type").String()
			}
			if item.Get("mem").Exists() {
				gpu.Memory = item.Get("mem").String()
			}

			if item.Get("memclock").Exists() {
				gpu.MemClock = item.Get("memclock").String()
			}
			if item.Get("coreclock").Exists() {
				gpu.CoreClock = item.Get("coreclock").String()
			}

			if item.Get("power").Exists() {
				gpu.Power = item.Get("power").String()
			}

			if item.Get("power").Exists() {
				gpu.Power = item.Get("power").String()
			}

			if item.Get("temp").Exists() {
				gpu.Temperature = item.Get("temp").String()
			}

			if item.Get("fanspeed").Exists() {
				gpu.FanSpeed = item.Get("fanspeed").String()
			}

			gpu.Driver = ""
			gpu.MaxPower = ""

			m[i] = gpu
			i += 1
		}

		return m, nil
	} else {
		return m, err
	}

	return m, errors.New("获取显卡信息失败")
}
func GetArgs(args string) []string {
	if len(args) > 0 {
		return strings.Split(args, ",")
	}

	var result []string
	result = make([]string, 10)
	return result
}

type CurrentMinerConfig struct {
	Name     string `json:"kernel"`  // 内核名称
	Version  string `json:"version"` // 版本
	Fullpath string `json:"path"`    // 全路径
	Command  string `json:"command"` // 启动程序名称
	Parames  string `json:"params"`  // 启动参数
}

func GetCurrentMiner() (*CurrentMinerConfig, error) {
	var m CurrentMinerConfig
	content, _ := ioutil.ReadFile("/home/ymminer/etc/miner.json")

	err := json.Unmarshal(content, &m)
	return &m, err
}
