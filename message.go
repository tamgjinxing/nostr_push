package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

// 生成指定长度的随机16进制的字符串
func GenerateRandomString(length int) string {
	// 定义字符集合
	charset := "ABCDEFabcdef0123456789"

	// 生成随机字符串
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}

	return string(randomString)
}

// 生成relay的订阅信息（req）
// reqId reqId
func GenerateSubscribeMsg(reqId string, filters map[string]interface{}) string {
	jsonMap := make(map[string]interface{})
	for key, value := range filters {
		jsonMap[key] = value
	}

	jsonArray := []interface{}{"REQ", reqId, jsonMap}

	jsonString, err := json.Marshal(jsonArray)
	if err != nil {
		return "" // 返回空字符串或者其他错误处理逻辑
	}

	return string(jsonString)
}

// 生成退订消息
func GenerateUnSubscribeMsg(reqId string) string {
	jsonArray := []interface{}{"CLOSE", reqId}

	jsonString, err := json.Marshal(jsonArray)
	if err != nil {
		return "" // 返回空字符串或者其他错误处理逻辑
	}

	return string(jsonString)
}

func GenerateAuthMsg(relay string, challenge string) (string, error) {
	tags := make([]Tag, 0, 2) // 初始长度为0，容量为5

	iTag := []string{"relay", relay}
	iTag2 := []string{"challenge", challenge}

	tags = append(tags, iTag)
	tags = append(tags, iTag2)

	event := Event{
		CreatedAt: time.Now().Unix(),
		Kind:      22242,
		Content:   "",
		Tags:      tags,
	}

	event.Sign(config.PushBotInfo.PrivateKey)

	return event.ToAuthString()
}
