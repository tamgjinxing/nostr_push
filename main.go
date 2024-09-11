package main

import (
	"os"
	"time"
)

func init() {
	eventChannel = make(chan HanlderEventInfo, 64)
	groupRelayChannel = make(chan HanlderEventInfo, 100)
	monentChannel = make(chan HanlderEventInfo, 64)
	inviteToGroupChannel = make(chan HanlderEventInfo, 64)

	// 创建 HTTP 客户端
	HttpClient = NewHTTPClient(HTTPClientConfig{
		Timeout:         10 * time.Second,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	})
}

func main() {
	env := os.Getenv("RUN_ENV")

	if env == "" {
		log.Println("环境变量 RUN_ENV 未设置或为空")
		return
	} else {
		log.Printf("环境变量 RUN_ENV 的值:%s\n", env)
	}

	relays := []MonitoringRelaysInfo{}
	configName := "config.json"
	relaysFileName := "relays.txt"

	isTest := false

	if env != "pro" {
		isTest = true
	}

	log.Printf("isTest:%t\n", isTest)
	if isTest {
		configName = "config_test.json"
		relaysFileName = "relays_test.txt"
	}

	// Initialize push management
	pushManager = NewPushManager(androidPusher, iosPusher)

	// Reading configuration files
	err := ReadConfig(configName)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	// For ios push
	go InitClientMap()

	ReadRelaysFile(relaysFileName)

	if isTest {
		relays = []MonitoringRelaysInfo{
			{RelayUrl: "ws://127.0.0.1:6970", SubKindsLimits: Oxchat_subscribeKinds, GroupRelayFlag: false},
			{RelayUrl: "ws://127.0.0.1:5577", SubKindsLimits: Oxchat_group_subscribeKinds, GroupRelayFlag: true},
		}
	} else {
		needMonitorRelays := GetMonitorChatRelays(config.TopInfo.TopN)
		for _, relay := range needMonitorRelays {
			relays = append(relays, MonitoringRelaysInfo{RelayUrl: relay, SubKindsLimits: Oxchat_subscribeKinds, GroupRelayFlag: false})
		}

		relays = append(relays, MonitoringRelaysInfo{RelayUrl: Group_Relay_0xchat, SubKindsLimits: Oxchat_group_subscribeKinds, GroupRelayFlag: false})
	}

	go ConnectToInitRelays(relays)

	go SendHeartbeat()

	go HandleEvent(eventChannel)

	go HandleGroupRelayEvent(groupRelayChannel)

	go MonmentEvent(monentChannel)

	go InviteToGroupEvent(inviteToGroupChannel)

	// Capture the interrupt signal and gracefully close the connection
	go HandleInterrupt()

	// Prevent the main function from exiting
	select {}
}
