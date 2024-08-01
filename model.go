package main

import "github.com/nbd-wtf/go-nostr"

type MemberDTO struct {
	ChannelId         string `json:"channelId"`
	MuteNotifications int    `json:"muteNotifications"`
	UserPubKey        string `json:"userPubKey"`
}

type ChannelInfoDTO struct {
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	About       string `json:"about"`
	Picture     string `json:"picture"`
	Owner       string `json:"owner"`
}

type UserInfoDTO struct {
	PublicKey string    `json:"publicKey"`
	DeviceId  string    `json:"deviceId"`
	Relays    []string  `json:"relays"`
	Kinds     []int     `json:"kinds"`
	Level     int       `json:"level"`
	ETags     nostr.Tag `json:"#e"`
	PTags     nostr.Tag `json:"#p"`
	Online    int       `json:"online"`
	Name      string    `json:"name"`
}

type UserInfo4Cache struct {
	Name   string `json:"name"`
	Pubkey string `json:"pubkey"`
}

type GroupInfo struct {
	GroupName string `json:"groupName"`
	GroupId   string `json:"groupId"`
	// GroupPic  string `json:"groupPic"`
	GroupType   string `json:"groupType"`
	GroupStatus string `json:"groupStatus"`
}
