package main

import (
	"errors"
	"fmt"
	"strings"
)

// Pusher 是推送接口，所有推送服务都要实现该接口
type Pusher interface {
	Push(message string, deviceID string, title string, isCallPush bool, groupId string) error
}

// AndroidPusher 实现安卓推送服务
type AndroidPusher struct {
}

func (a *AndroidPusher) Push(message string, deviceID string, title string, isCallPush bool, groupId string) error {
	// 这里实现安卓推送的具体逻辑
	PushAndroid(deviceID, title, message, isCallPush, groupId)
	return nil
}

// IOSPusher 实现iOS推送服务
type IOSPusher struct{}

func (i *IOSPusher) Push(message string, deviceID string, title string, isCallPush bool, groupId string) error {
	// 这里实现iOS推送的具体逻辑
	PushIos(deviceID, title, message)
	return nil
}

// PusherType 是推送服务类型的枚举
type PusherType int

const (
	Android PusherType = iota
	IOS
)

// PusherFactory 是推送服务工厂
type PusherFactory struct{}

func (f *PusherFactory) CreatePusher(t PusherType) (Pusher, error) {
	switch t {
	case Android:
		return &AndroidPusher{}, nil
	case IOS:
		return &IOSPusher{}, nil
	default:
		return nil, errors.New("unknown pusher type")
	}
}

// PushManager 是推送管理器，负责选择合适的推送服务
type PushManager struct {
	androidPusher Pusher
	iosPusher     Pusher
}

// NewPushManager 创建一个新的 PushManager
func NewPushManager(androidPusher, iosPusher Pusher) *PushManager {
	return &PushManager{
		androidPusher: androidPusher,
		iosPusher:     iosPusher,
	}
}

// PushMessage 根据设备号选择合适的推送服务发送消息
func (m *PushManager) PushMessage(eventIdAndPubkey string, message string, deviceID string, title string, isCallPush bool, groupId string) {
	var pusher Pusher
	if strings.HasPrefix(deviceID, "com") {
		pusher = m.iosPusher
	} else {
		pusher = m.androidPusher
	}

	isPushed, _ := RecentKeysExist(pushToPubkeyEidRedisKey + eventIdAndPubkey)

	if !isPushed {
		if err := pusher.Push(message, deviceID, title, isCallPush, groupId); err != nil {
			fmt.Printf("Failed to push message to device %s: %v\n", deviceID, err)
		} else {
			SaveRecentKey(pushToPubkeyEidRedisKey + eventIdAndPubkey)
		}
	}
}
