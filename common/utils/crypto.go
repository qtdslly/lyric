package utils

import (
	"crypto/aes"
	"crypto/cipher"
	cr "crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"strings"

)

var key = "0g.DockWallet*)5"

func GetMinerConfigContent() string {
	data, err := ioutil.ReadFile("/opt/miner/miner-config.json")
	if err != nil {
		return ""
	}

	return string(data)
}

func GetNfsnetConfigContent() string {
	data, err := ioutil.ReadFile("/opt/miner/nfsnet-config.json")
	if err != nil {
		return ""
	}

	return string(data)
}

func Decrypt(content string) string {
	hex_data, _ := hex.DecodeString(content)
	decrypted := AesDecryptCFB(hex_data, []byte(key))
	return string(decrypted)
}

func AesDecryptCFB(encrypted []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		return []byte{}
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
func Encrypt(content string) string {
	encrypted := AesEncryptCFB([]byte(content), []byte(key))
	if len(encrypted) == 0 {
		return ""
	}
	result := hex.EncodeToString(encrypted)
	return string(result)
}
func AesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(cr.Reader, iv); err != nil {
		return
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}

func GetSoftVersion(softPath string) string {
	data, err := BashCommand(softPath+" -v", 5)
	if err != nil {
		return ""
	}
	data = strings.Trim(data, "\n")
	return data
}
