package main

import (
	"encoding/json"
	"os"
)

type Config struct {
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

	TopInfo struct {
		TopN int `json:"topN"`
	}
}

var config Config

func ReadConfig(filename string) error {
	log.Printf("load config file:%s\n", filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}
