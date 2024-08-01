package main

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	log                    = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Caller().Logger()
	PRIVATE_MSG_PUSH_KINDS = []int{4, 1059, 44}

	MONENT_REPLY_PUSH_KINDS = []int{1}
	MONEMT_LIKE_PUSH_KINDS  = []int{7}

	clients []*WebSocketClient
)

const (
	PrivateMsgPushKind = 1059

	CallPushKind = "25050"
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
