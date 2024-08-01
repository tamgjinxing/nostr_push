package main

import (
	"github.com/nbd-wtf/go-nostr"
	"github.com/sideshow/apns2"
)

var (
	clientMap = make(map[string]*apns2.Client)

	// 创建推送服务工厂
	factory = &PusherFactory{}

	// 创建安卓和iOS推送服务
	androidPusher, _ = factory.CreatePusher(Android)
	iosPusher, _     = factory.CreatePusher(IOS)

	// 创建推送管理器
	pushManager *PushManager

	eventChannel         chan HanlderEventInfo
	groupRelayChannel    chan HanlderEventInfo
	monentChannel        chan HanlderEventInfo
	inviteToGroupChannel chan HanlderEventInfo

	GroupRelayChannelKinds = []int{
		nostr.KindSimpleGroupChatMessage,
		nostr.KindSimpleGroupThread,
		nostr.KindSimpleGroupReply,
		nostr.KindSimpleGroupMetadata,
		nostr.KindSimpleGroupMembers}

	MonmentChannelKinds = []int{
		nostr.KindTextNote,
		nostr.KindReaction}
)

func init() {
	eventChannel = make(chan HanlderEventInfo, 64)
	groupRelayChannel = make(chan HanlderEventInfo, 100)
	monentChannel = make(chan HanlderEventInfo, 64)
	inviteToGroupChannel = make(chan HanlderEventInfo, 64)
}

func main() {

	// 初始化推送管理
	pushManager = NewPushManager(androidPusher, iosPusher)

	// 读取配置文件
	err := ReadConfig("config.json")
	if err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	// 用于ios推送的
	InitClientMap()

	oxchat_group_subscribeKinds := []SubKindsLimit{
		{SubscribeKinds: []int{9, 11, 12, 9000}, Limit: 0},
		{SubscribeKinds: []int{39000, 39002}, Limit: 1},
	}

	oxchat_subscribeKinds := []SubKindsLimit{
		{SubscribeKinds: []int{1059, 42}, Limit: 0},
	}

	//连接到各个relay服务
	relays := []MonitoringRelaysInfo{
		// {RelayUrl: "wss://relay.0xchat.com", SubscribeKinds: []int{42, 1059}, isGroupRelay: false},
		// {RelayUrl: "wss://groups.0xchat.com", SubscribeKinds: []int{39000, 39002}},
		{RelayUrl: "ws://127.0.0.1:5578", SubKindsLimits: oxchat_subscribeKinds, GroupRelayFlag: false},
		// {RelayUrl: "ws://127.0.0.1:5577", SubscribeKinds: []int{39000, 39002}, Limit: 1},
		{RelayUrl: "ws://127.0.0.1:5577", SubKindsLimits: oxchat_group_subscribeKinds, GroupRelayFlag: true},
		// Add more instances as needed
	}

	go ConnectToInitRelays(relays)

	go HandleEvent(eventChannel)

	go HandleGroupRelayEvent(groupRelayChannel)

	go MonmentEvent(monentChannel)

	go InviteToGroupEvent(inviteToGroupChannel)

	// 捕获中断信号并优雅地关闭连接
	go HandleInterrupt(clients)

	// Prevent the main function from exiting
	select {}
}
