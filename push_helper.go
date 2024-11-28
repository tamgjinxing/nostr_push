package main

import (
	"errors"
	"fmt"
	"strings"
)

// Pusher is the push interface that all push services must implement
type Pusher interface {
	Push(message string, deviceID string, title string, isCallPush bool, groupId string) error
}

// AndroidPusher Implement Android push service
type AndroidPusher struct {
}

func (a *AndroidPusher) Push(message string, deviceID string, title string, isCallPush bool, groupId string) error {
	PushAndroid(deviceID, title, message, isCallPush, groupId)
	return nil
}

type IOSPusher struct{}

func (i *IOSPusher) Push(message string, deviceID string, title string, isCallPush bool, groupId string) error {
	PushIos(deviceID, title, message, isCallPush, groupId)
	return nil
}

type PusherType int

const (
	Android PusherType = iota
	IOS
)

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

type PushManager struct {
	androidPusher Pusher
	iosPusher     Pusher
}

func NewPushManager(androidPusher, iosPusher Pusher) *PushManager {
	return &PushManager{
		androidPusher: androidPusher,
		iosPusher:     iosPusher,
	}
}

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
