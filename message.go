package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

// Generates a random hexadecimal string of the specified length
func GenerateRandomString(length int) string {
	charset := "ABCDEFabcdef0123456789"

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}

	return string(randomString)
}

// reqId reqId
func GenerateSubscribeMsg(reqId string, filters map[string]interface{}) string {
	jsonMap := make(map[string]interface{})
	for key, value := range filters {
		jsonMap[key] = value
	}

	jsonArray := []interface{}{"REQ", reqId, jsonMap}

	jsonString, err := json.Marshal(jsonArray)
	if err != nil {
		return ""
	}

	return string(jsonString)
}

// Generate an unsubscribe message
func GenerateUnSubscribeMsg(reqId string) string {
	jsonArray := []interface{}{"CLOSE", reqId}

	jsonString, err := json.Marshal(jsonArray)
	if err != nil {
		return ""
	}

	return string(jsonString)
}

func GenerateAuthMsg(relay string, challenge string) (string, error) {
	tags := make([]Tag, 0, 2)

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
