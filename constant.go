package main

import (
	"os"
	"sync"

	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog"
	"github.com/sideshow/apns2"
)

var (
	log                    = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}).With().Timestamp().Logger()
	PRIVATE_MSG_PUSH_KINDS = []int{4, 1059, 44}

	MONENT_REPLY_PUSH_KINDS = []int{1}
	MONEMT_LIKE_PUSH_KINDS  = []int{7}

	clients []*WebSocketClient

	clientsMap = make(map[string]*WebSocketClient)

	mu sync.RWMutex

	needMonitorRelays []string

	// group relay:
	Oxchat_group_subscribeKinds = []SubKindsLimit{
		{SubscribeKinds: []int{9, 11, 12, 9000}, Limit: 0},
		{SubscribeKinds: []int{39000, 39002}, Limit: 1},
	}

	// normal relay:
	Oxchat_subscribeKinds = []SubKindsLimit{
		{SubscribeKinds: []int{1059, 42, 7, 1, 9735}, Limit: 0},
	}
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

	HttpClient *HTTPClient
)

const (
	PrivateMsgPushKind = 1059

	CallPushKind = "25050"
)

const (
	Group_Relay_0xchat = "wss://groups.0xchat.com"
)

const (
	MembersRedisKey    = "channel:members:"
	GroupInfoKRedisKey = "channel:info:"

	UserPushInfoRedisKey      = "user:push:info:"
	UserInfoRedisKey          = "user:info:"
	groupRelayGroupIdRedisKey = "groupRelayGroupId:"
	groupRelayMembersRedisKey = "groupRelayMembers:"

	pushToPubkeyEidRedisKey = "pushToPubkeyEid:"
)

const (
	Default_push_title  = "0xchat"
	Default_private_msg = "Received a private message"
	Default_channel_msg = ""
	Default_call_msg    = "Received a call request"
	Default_reply_msg   = "%s reply a note your were mentioned in"
	Default_like_msg    = "%s like a note your were mentioned in"
)
