package main

import (
	"encoding/json"
	"os"
)

// Config 结构体用于存储配置信息
type Config struct {
	RelayInfo struct {
		DatabasePath     string `json:"databasePath"`
		Port             string `json:"port"`
		Domain           string `json:"domain"`
		RelayName        string `json:"name"`
		RelayDescription string `json:"description"`
		RelayContact     string `json:"contact"`
		RelayIcon        string `json:"icon"`
		RelayPubkey      string `json:"pubkey"`
	} `json:"RelayInfo"`

	IOSPushInfo struct {
		BundleIds string `json:"bundleIds"`
		P12Pathes string `json:"p12Pathes"`
		Passwords string `json:"passwords"`
	} `json:"IosPushInfo"`

	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"Redis"`

	PushBotInfo struct {
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
	}
}

var config Config // 全局变量用于存储配置信息

// readConfig 函数用于读取配置文件并将其内容存储在全局变量中
func ReadConfig(filename string) error {
	log.Printf("加载配置文件:%s\n", filename)
	// 读取配置文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// 解析配置文件内容到 Config 结构体
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}
