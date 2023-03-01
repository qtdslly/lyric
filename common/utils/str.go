package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetUserWallet(data string) string {
	parttern := `0x[0-9a-zA-Z]{40}`
	re, err := regexp.Compile(parttern)
	if err != nil {
		return ""
	}

	wallet := string(re.Find([]byte(data)))
	return wallet
}

func GetFileLine(filename string) int64 {
	data, err := BashCommand("wc -l "+filename+" | awk '{print $1}'", 5)
	if err != nil {
		return 0
	}

	data = strings.Replace(data, "\n", "", -1)

	line, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0
	}
	return line
}
