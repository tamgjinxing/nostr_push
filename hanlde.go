package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

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
					go ChannelPush(ev.Event)
				case PrivateMsgPushKind:
					go PrivatePush(ev.Event)
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
					go PublicGroupPush(ev.Event)
				case nostr.KindSimpleGroupThread:
					go PublicGroupPush(ev.Event)
				case nostr.KindSimpleGroupReply:
					go PublicGroupPush(ev.Event)
				case nostr.KindSimpleGroupMetadata:
					go SaveGroupRelayGroupInfo(ev.Event)
				case nostr.KindSimpleGroupMembers:
					go SaveGroupRelayGroupMembers(ev.Event)
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
					go MonmentPush(ev.Event)
				case nostr.KindReaction:
					go MonmentPush(ev.Event)
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
					go InviteToGroupHandler(ev)
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
	groupType := ""
	groupStatus := ""

	for _, tag := range tags {
		if len(tag) == 2 {
			tag0 := tag[0]
			tag1 := tag[1]

			if tag0 == "d" {
				groupId = tag1
			}

			if tag0 == "name" {
				groupName = tag1
			}
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
		GroupId:     groupId,
		GroupName:   groupName,
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
		for _, value := range members {
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
				sort.Ints(userInfo.Kinds)
				_, found := slices.BinarySearch(userInfo.Kinds, nostr.KindChannelMessage)
				groupFound := StringInSlice(userInfo.ETags, groupId)
				if found && groupFound {
					if channelInfo != nil {
						go SendMessageToMember(userInfo, channelInfo.ChannelName, event.Content, event.ID, sendName)
					} else {
						go SendMessageToMember(userInfo, "", event.Content, event.ID, sendName)
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
		log.Println("Kind=1059 and K tag value isn't 25050, no need push!!!")
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

				go pushManager.PushMessage(event.ID+"_"+userInfo.PublicKey, defaultMsg, userInfo.DeviceId, Default_push_title, isCallPush, "")
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
			log.Println(member)
			if member == event.PubKey {
				log.Printf("The event is send by self,is no need push. pubkey=:%s", member)
				continue
			}

			userInfo := GetUserInfoFromRedis(member)

			if userInfo != nil {
				sort.Ints(userInfo.Kinds)
				_, found := slices.BinarySearch(userInfo.Kinds, event.Kind)
				groupFound := StringInSlice(userInfo.ETags, groupId)
				if found && groupFound {
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

						go pushManager.PushMessage(event.ID+"_"+userInfo.PublicKey, message, userInfo.DeviceId, title, false, groupId)
					}
				}
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
	} else if event.Kind == nostr.KindReaction {
		sendMsg = defaultLikeMsg
	}

	toPubkeys = RemoveDuplicates(toPubkeys)
	for _, toPubkey := range toPubkeys {
		if toPubkey == event.PubKey {
			log.Printf("The person being replied to is the same as the person replying, no need to push a notification.")
			continue
		}

		log.Printf("The person being replied to is the user who needs to receive a notification. pubkey:%s\n", toPubkey)
		userInfo := GetUserInfoFromRedis(toPubkey)
		if userInfo != nil {
			sort.Ints(userInfo.Kinds)
			_, match := slices.BinarySearch(userInfo.Kinds, event.Kind)
			if match {
				if userInfo.Online == 0 && userInfo.DeviceId != "" {
					go pushManager.PushMessage(event.ID+"_"+userInfo.PublicKey, sendMsg, userInfo.DeviceId, "0xchat", false, "")
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
					log.Println("Upon receiving an invitation to join the group, send a request to get basic group information and subscribe to retrieve the group details and member list!")
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

func SendMessageToMember(userInfo *UserInfoDTO, groupName string, message string, eventId string, sendName string) {
	if userInfo.Online == 0 && userInfo.DeviceId != "" {

		if sendName != "" {
			message = sendName + ":" + message
		}

		title := groupName

		if title == "" {
			title = "0xchat"
		}

		go pushManager.PushMessage(eventId+"_"+userInfo.PublicKey, message, userInfo.DeviceId, title, false, "")
	}
}

func StringInSlice(list []string, a string) bool {
	for _, b := range list {
		if strings.Contains(b, a) {
			return true
		}
	}
	return false
}
