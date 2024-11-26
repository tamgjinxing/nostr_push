package main

import (
	"slices"
	"sort"
	"sync"

	"github.com/nbd-wtf/go-nostr"
)

// 定义一个 Goroutine 池
var workerPool = make(chan struct{}, 20) // 允许20个并发任务
var workerPoolMore = make(chan struct{}, 100)

var workerPool250 = make(chan struct{}, 500)

func ChannelPush2(event *nostr.Event) {
	tags := event.Tags
	var groupId string

	if len(tags) > 0 {
		for _, tag := range tags {
			if len(tag) > 2 && tag[0] == "e" {
				groupId = tag[1]
				break // 找到第一个 "e" 标签即可，提前结束循环
			}
		}
	}

	if groupId != "" {
		var wg sync.WaitGroup
		var channelInfo *ChannelInfoDTO
		var members map[string]MemberDTO
		var userInfo4Cache *UserInfo4Cache

		wg.Add(1)
		go func() {
			defer wg.Done()
			channelInfo, _ = GetChannelInfoFromRedis(groupId)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			members, _ = GetMembersFromRedis(groupId)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			userInfo4Cache, _ = GetUserInfo4Cache(event.PubKey)
		}()

		wg.Wait()

		log.Printf("members:%d\n", len(members))
		for _, value := range members {
			if event.PubKey == value.UserPubKey {
				continue // 跳过消息发送者
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
					// 使用 Goroutine 池控制并发数量
					if len(members) > 1000 {
						workerPoolMore <- struct{}{} // 占用池中的一个 Goroutine
						go func(userInfo *UserInfoDTO) {
							defer func() { <-workerPoolMore }() // 任务完成后释放 Goroutine
							if channelInfo != nil {
								SendMessageToMember(userInfo, channelInfo.ChannelName, event.Content, event.ID, sendName)
							} else {
								SendMessageToMember(userInfo, "", event.Content, event.ID, sendName)
							}
						}(userInfo)
					} else {
						workerPool <- struct{}{} // 占用池中的一个 Goroutine
						go func(userInfo *UserInfoDTO) {
							defer func() { <-workerPool }() // 任务完成后释放 Goroutine
							if channelInfo != nil {
								SendMessageToMember(userInfo, channelInfo.ChannelName, event.Content, event.ID, sendName)
							} else {
								SendMessageToMember(userInfo, "", event.Content, event.ID, sendName)
							}
						}(userInfo)
					}
				}
			}
		}
	}
}

func ChannelPush3(event *nostr.Event) {
	tags := event.Tags
	var groupId string

	// 查找 groupId
	if len(tags) > 0 {
		for _, tag := range tags {
			if len(tag) > 2 && tag[0] == "e" {
				groupId = tag[1]
				break
			}
		}
	}

	// 如果找到 groupId，则进行处理
	if groupId != "" {
		// 异步获取 ChannelInfo 和成员信息
		var wg sync.WaitGroup
		var channelInfo *ChannelInfoDTO
		var members map[string]MemberDTO

		// 异步获取群组信息
		wg.Add(1)
		go func() {
			defer wg.Done()
			channelInfo, _ = GetChannelInfoFromRedis(groupId)
		}()

		// 异步获取成员列表
		wg.Add(1)
		go func() {
			defer wg.Done()
			members, _ = GetMembersFromRedis(groupId)
		}()

		// 等待所有数据获取完成
		wg.Wait()

		log.Printf("members:%d\n", len(members))

		// 异步处理每个成员的消息推送
		for _, member := range members {
			// 消息发送者无需推送
			if event.PubKey == member.UserPubKey {
				continue
			}

			workerPool250 <- struct{}{} // 限制并发数，防止过载
			wg.Add(1)
			go func(member *MemberDTO) {
				defer func() {
					<-workerPool250 // 释放 Goroutine
					wg.Done()       // 标记任务完成
				}()

				// 获取发送者名称（异步获取缓存用户信息）
				sendName := ""
				userInfo4Cache, _ := GetUserInfo4Cache(event.PubKey)
				if userInfo4Cache != nil {
					sendName = userInfo4Cache.Name
				}

				// 异步获取目标用户信息
				userInfo := GetUserInfoFromRedis(member.UserPubKey)

				if userInfo != nil {
					// 判断用户的 Kind 和群组订阅情况
					sort.Ints(userInfo.Kinds)
					_, found := slices.BinarySearch(userInfo.Kinds, nostr.KindChannelMessage)
					groupFound := StringInSlice(userInfo.ETags, groupId)

					// 满足条件才进行推送
					if found && groupFound {
						if channelInfo != nil {
							SendMessageToMember(userInfo, channelInfo.ChannelName, event.Content, event.ID, sendName)
						} else {
							SendMessageToMember(userInfo, "", event.Content, event.ID, sendName)
						}
					}
				}
			}(&member)
		}

		// 等待所有成员处理完成
		wg.Wait()
	}
}
