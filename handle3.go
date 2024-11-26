package main

import (
	"strconv"
	"sync"

	"github.com/nbd-wtf/go-nostr"
)

// ChannelPush 使用 MemberDTO 并发推送消息
func ChannelPush4(event *nostr.Event) {
	tags := event.Tags
	var groupId string

	// 提取 groupId
	if len(tags) > 0 {
		for _, tag := range tags {
			if len(tag) > 2 && tag[0] == "e" {
				groupId = tag[1]
			}
		}
	}

	// 如果有 groupId，则继续处理
	if groupId != "" {
		var wg sync.WaitGroup
		var channelInfo *ChannelInfoDTO
		members := make(map[string]MemberDTO)
		// var userInfo4Cache *UserInfo4Cache

		wg.Add(1)
		go func() {
			defer wg.Done()
			channelInfo, _ = GetChannelInfoFromRedis(groupId)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			// members, _ = GetMembersFromRedis(groupId)
			for i := 0; i < 1000; i++ {
				members["pubkey"+strconv.Itoa(i)] = MemberDTO{UserPubKey: "pubkey" + strconv.Itoa(i), ChannelId: "14f71ee3b7c8af6206746ecaa1ecd7fd5cb6edbbf447ca5ab97be24dc1a70078"}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			// userInfo4Cache, _ = GetUserInfo4Cache(event.PubKey)
		}()

		wg.Wait()

		// 创建一个 WaitGroup 以等待所有并发的消息推送完成

		// 控制并发的 workers 数量
		workers := 200
		// 创建一个用于存储待处理任务的 channel
		jobs := make(chan *MemberDTO, 1000)

		// 启动 workers goroutines
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for member := range jobs {
					// 如果发送者与接收者是同一个人，则跳过
					if event.PubKey == member.UserPubKey {
						continue
					}

					sendName := ""
					// if userInfo4Cache != nil {
					// 	sendName = userInfo4Cache.Name
					// }

					userInfo := UserInfoDTO{PublicKey: member.UserPubKey}

					SendMessageToMember(&userInfo, channelInfo.ChannelName, event.Content, event.ID, sendName)
					// userInfo := GetUserInfoFromRedis(member.UserPubKey)

					// if userInfo != nil {
					// 	// 检查用户是否订阅了该频道
					// 	sort.Ints(userInfo.Kinds)
					// 	_, found := slices.BinarySearch(userInfo.Kinds, nostr.KindChannelMessage)
					// 	groupFound := StringInSlice(userInfo.ETags, groupId)
					// 	if found && groupFound {
					// 		// 如果频道信息存在，发送带频道名称的消息
					// 		if channelInfo != nil {
					// 			SendMessageToMember(userInfo, channelInfo.ChannelName, event.Content, event.ID, sendName)
					// 		} else {
					// 			// 如果频道信息不存在，发送不带频道名称的消息
					// 			SendMessageToMember(userInfo, "", event.Content, event.ID, sendName)
					// 		}
					// 	}
					// }
				}
			}()
		}

		// 将每个成员推送到 jobs channel 中
		for _, member := range members {
			jobs <- &member
		}
		close(jobs)

		// 等待所有消息推送完成
		wg.Wait()
	}
}
