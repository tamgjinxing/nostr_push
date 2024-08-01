package main

import (
	"fmt"
	"slices"
	"sort"
	"sync"

	"github.com/nbd-wtf/go-nostr"
)

// handle the req return event
func HandleEvent(events chan HanlderEventInfo) {
	for {
		ev, ok := <-events
		if ok {
			if ev.Event != nil {
				switch ev.Event.Kind {
				case nostr.KindChannelMessage:
					ChannelPush(ev.Event)
				case PrivateMsgPushKind:
					PrivatePush(ev.Event)
				}
			}
		}
	}
}

func HandleGroupRelayEvent(events chan HanlderEventInfo) {
	for {
		ev, ok := <-events
		if ok {
			if ev.Event != nil {
				switch ev.Event.Kind {
				case nostr.KindSimpleGroupChatMessage:
					PublicGroupPush(ev.Event)
				case nostr.KindSimpleGroupThread:
					PublicGroupPush(ev.Event)
				case nostr.KindSimpleGroupReply:
					PublicGroupPush(ev.Event)
				case nostr.KindSimpleGroupMetadata:
					SaveGroupRelayGroupInfo(ev.Event)
				case nostr.KindSimpleGroupMembers:
					SaveGroupRelayGroupMembers(ev.Event)
				}
			}
		}
	}
}

func MonmentEvent(events chan HanlderEventInfo) {
	for {
		ev, ok := <-events
		if ok {
			if ev.Event != nil {
				switch ev.Event.Kind {
				case nostr.KindTextNote:
					MonmentPush(ev.Event)
				case nostr.KindReaction:
					MonmentPush(ev.Event)
				}
			}
		}
	}
}

func InviteToGroupEvent(events chan HanlderEventInfo) {
	for {
		ev, ok := <-events
		if ok {
			if ev.Event != nil {
				switch ev.Event.Kind {
				case 9000:
					InviteToGroupHandler(ev)
				}
			}
		}
	}
}

func SaveGroupRelayGroupMembers(event *nostr.Event) {
	tags := event.Tags

	groupId := ""
	members := []string{}

	for _, tag := range tags {

		if len(tag) == 2 {
			tag0 := tag[0]
			tag1 := tag[1]

			// log.Printf("tag0:%s", tag[0])
			// log.Printf("tag1:%s\n", tag[1])

			if tag0 == "d" {
				groupId = tag1
			}

			if tag0 == "p" {
				members = append(members, tag1)
			}
		}
	}

	if groupId != "" {
		go PutGroupRelayGroupMembersToRedis(groupId, members)
	}
}

func SaveGroupRelayGroupInfo(event *nostr.Event) {
	tags := event.Tags

	groupId := ""
	groupName := ""
	// groupPic := ""
	groupType := ""
	groupStatus := ""

	for _, tag := range tags {
		if len(tag) == 2 {
			tag0 := tag[0]
			tag1 := tag[1]

			// log.Printf("tag0:%s", tag[0])
			// log.Printf("tag1:%s\n", tag[1])

			if tag0 == "d" {
				groupId = tag1
			}

			if tag0 == "name" {
				groupName = tag1
			}

			// if tag0 == "picture" {
			// 	groupPic = tag1
			// }
		}

		if len(tag) == 1 {
			if tag[0] == "public" || tag[0] == "private" {
				groupType = tag[0]
			}

			if tag[0] == "open" || tag[0] == "closed" {
				groupStatus = tag[0]
			}
		}
	}

	groupInfo := GroupInfo{
		GroupId:   groupId,
		GroupName: groupName,
		// GroupPic:  groupPic,
		GroupType:   groupType,
		GroupStatus: groupStatus,
	}

	if groupId != "" {
		go PutGroupRelayGroupInfoToRedis(groupInfo)
	}
}

func ChannelPush(event *nostr.Event) {
	tags := event.Tags
	var groupId string

	if len(tags) > 0 {
		for _, tag := range tags {
			if len(tag) > 2 {
				// log.Println(tag)

				if tag[0] == "e" {
					groupId = tag[1]
				}
			}
		}
	}

	if groupId != "" {
		channelInfo, _ := GetChannelInfoFromRedis(groupId)
		members, _ := GetMembersFromRedis(groupId)

		userInfo4Cache, _ := GetUserInfo4Cache(event.PubKey)
		var wg sync.WaitGroup
		for _, value := range members {
			wg.Add(1)
			// message sender is no need push
			if event.PubKey == value.UserPubKey {
				continue
			}

			sendName := ""
			if userInfo4Cache != nil {
				sendName = userInfo4Cache.Name
			}

			userInfo := GetUserInfoFromRedis(value.UserPubKey)

			if userInfo != nil {
				// 对切片进行排序
				sort.Ints(userInfo.Kinds)
				_, found := slices.BinarySearch(userInfo.Kinds, nostr.KindChannelMessage)
				log.Printf("found:%t\n", found)
				if found {
					if channelInfo != nil {
						go SendMessageToMember(userInfo, channelInfo.ChannelName, event.Content, &wg, event.ID, sendName)
					} else {
						go SendMessageToMember(userInfo, "", event.Content, &wg, event.ID, sendName)
					}
				}
			}
		}
	}
}

func PrivatePush(event *nostr.Event) {
	tags := event.Tags
	toPubkey := ""
	isCallPush := false
	isNeedPush := true

	defaultMsg := Default_private_msg

	if len(tags) > 0 {
		hasKTag := false

		for _, tag := range tags {
			if len(tag) >= 2 {
				if tag[0] == "p" {
					toPubkey = tag[1]
				}

				if event.Kind == PrivateMsgPushKind {
					if tag[0] == "k" {
						if tag[1] == CallPushKind {
							isCallPush = true
						} else {
							isNeedPush = false
						}

						hasKTag = true
					}
				}
			}
		}

		if !hasKTag {
			isNeedPush = true
		}
	}

	if !isNeedPush {
		log.Println("1059的非25050消息,不进行推送!!!")
		return
	}

	userInfo := GetUserInfoFromRedis(toPubkey)
	if userInfo != nil {
		match := AnyMatch(PRIVATE_MSG_PUSH_KINDS, userInfo.Kinds)
		if match {
			if userInfo.Online == 0 && userInfo.DeviceId != "" {
				if isCallPush {
					defaultMsg = Default_call_msg
				}

				pushManager.PushMessage(event.ID+"_"+userInfo.PublicKey, defaultMsg, userInfo.DeviceId, Default_push_title, isCallPush, "")
			}
		}
	}
}

func PublicGroupPush(event *nostr.Event) {
	tags := event.Tags
	groupId := ""

	for _, tag := range tags {
		if len(tag) >= 2 {
			if tag[0] == "h" {
				groupId = tag[1]
				break
			}
		}
	}

	if groupId != "" {
		groupInfo := GetGroupRelayGroupInfoFromRedis(groupId)
		members, err := getGroupRelayMembersFromRedis(groupId)
		if err != nil {
			return
		}

		for _, member := range members {
			if member == event.PubKey {
				log.Printf("发送人自己不需要推送，pubkey=:%s", member)
				continue
			}

			userInfo := GetUserInfoFromRedis(member)

			if userInfo.Online == 0 && userInfo.DeviceId != "" {
				message := event.Content
				userInfo4Cache, _ := GetUserInfo4Cache(event.PubKey)
				if userInfo4Cache != nil {
					if userInfo4Cache.Name != "" {
						message = userInfo4Cache.Name + ":" + message
					}
				}

				title := groupInfo.GroupName
				if title == "" {
					title = Default_push_title
				}

				pushManager.PushMessage(event.ID+"_"+userInfo.PublicKey, message, userInfo.DeviceId, title, false, groupId)
			}
		}
	}
}

func MonmentPush(event *nostr.Event) {
	tags := event.Tags
	var toPubkeys []string

	sendMsg := ""
	defaultReplyMsg := fmt.Sprintf(Default_reply_msg, "someone")
	defaultLikeMsg := fmt.Sprintf(Default_like_msg, "someone")

	var matchKinds []int

	if len(tags) > 0 {
		for _, tag := range tags {
			if len(tag) >= 2 {
				if tag[0] == "p" {
					toPubkeys = append(toPubkeys, tag[1])
				}
			}
		}
	}

	sendUserInfo, _ := GetUserInfo4Cache(event.PubKey)

	if sendUserInfo != nil && sendUserInfo.Name != "" {
		defaultReplyMsg = fmt.Sprintf(Default_reply_msg, sendUserInfo.Name)
		defaultLikeMsg = fmt.Sprintf(Default_like_msg, sendUserInfo.Name)
	}

	if event.Kind == nostr.KindTextNote {
		sendMsg = defaultReplyMsg
		matchKinds = MONENT_REPLY_PUSH_KINDS
	} else if event.Kind == nostr.KindReaction {
		sendMsg = defaultLikeMsg
		matchKinds = MONEMT_LIKE_PUSH_KINDS
	}

	toPubkeys = RemoveDuplicates(toPubkeys)
	for _, toPubkey := range toPubkeys {
		if toPubkey == event.PubKey {
			log.Printf("被回复的人恰好是回复的人，不需要推送")
			continue
		}

		log.Printf("被回复的人，需要推送的用户pubkey:%s\n", toPubkey)
		userInfo := GetUserInfoFromRedis(toPubkey)
		if userInfo != nil {
			match := AnyMatch(matchKinds, userInfo.Kinds)
			if match {
				if userInfo.Online == 0 && userInfo.DeviceId != "" {
					pushManager.PushMessage(event.ID+"_"+userInfo.PublicKey, sendMsg, userInfo.DeviceId, "0xchat", false, "")
				}
			}
		}
	}
}

func InviteToGroupHandler(handleEventInfo HanlderEventInfo) {
	if handleEventInfo.Event != nil {
		event := handleEventInfo.Event

		tags := event.Tags
		invitedPubKeys := []string{}
		groupId := ""

		for _, tag := range tags {
			if len(tag) >= 2 {
				if tag[0] == "p" {
					invitedPubKeys = append(invitedPubKeys, tag[1])
				}

				if tag[0] == "h" {
					groupId = tag[1]
				}
			}
		}
		log.Println(invitedPubKeys)
		for _, invitedPubKey := range invitedPubKeys {
			if invitedPubKey == config.PushBotInfo.PublicKey {
				client := handleEventInfo.Client
				if client != nil {
					subId := GenerateRandomString(32)
					filters := map[string]interface{}{
						"kinds": []int{39000, 39002},
						"#d":    []string{groupId},
					}
					log.Println("接收到邀请入群通知，发送获取群基本信息和群成员订阅获取群信息和成员列表！")
					reqMsg := GenerateSubscribeMsg(subId, filters)
					err := client.SendMessage(reqMsg)
					if err != nil {
						log.Printf("Failed to send message to %s: %v", client.MonitorgRelaysInfo.RelayUrl, err)
					} else {
						log.Printf("Sent message to %s: %s\n", client.MonitorgRelaysInfo.RelayUrl, reqMsg)
					}
				}
			}
		}
	}
}

// SendMessageToMember 模拟发送消息给单个成员的函数
func SendMessageToMember(userInfo *UserInfoDTO, groupName string, message string, wg *sync.WaitGroup, eventId string, sendName string) {
	defer wg.Done()
	if userInfo.Online == 0 && userInfo.DeviceId != "" {

		if sendName != "" {
			message = sendName + ":" + message
		}

		title := groupName

		if title == "" {
			title = "0xchat"
		}

		pushManager.PushMessage(eventId+"_"+userInfo.PublicKey, message, userInfo.DeviceId, title, false, "")
	}
}
