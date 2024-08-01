package main

import "encoding/json"

func GetChannelInfoFromRedis(groupId string) (*ChannelInfoDTO, error) {
	value, err := GetRedisByKey(GroupInfoKRedisKey + groupId)
	if err != nil {
		return nil, err
	}

	var channelInfoDTO ChannelInfoDTO
	err = json.Unmarshal([]byte(value), &channelInfoDTO)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
	}

	return &channelInfoDTO, nil
}

func GetMembersFromRedis(groupId string) (map[string]MemberDTO, error) {
	var mapCache map[string]MemberDTO

	err := HGetJSON(MembersRedisKey, groupId, &mapCache)
	if err != nil {
		log.Printf("Error getting hash field: %v", err)
		return nil, err
	}
	return mapCache, nil
}

func GetUserInfoFromRedis(pubkey string) *UserInfoDTO {
	var userInfo UserInfoDTO
	_ = GetJSON(UserPushInfoRedisKey+pubkey, &userInfo)

	return &userInfo
}

// GetUserInfo4Cache 获取用户信息
func GetUserInfo4Cache(pubkey string) (*UserInfo4Cache, error) {
	var userInfo UserInfo4Cache
	err := HGetJSON(UserInfoRedisKey, pubkey, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func PutGroupRelayGroupInfoToRedis(groupInfo GroupInfo) error {
	err := PutJSON(groupRelayGroupIdRedisKey+groupInfo.GroupId, groupInfo)

	if err != nil {
		return err
	}

	return nil
}

func PutGroupRelayGroupMembersToRedis(groupId string, members []string) error {
	err := PutHashList(groupRelayMembersRedisKey, groupId, members)

	if err != nil {
		return err
	}

	return nil
}

func GetGroupRelayGroupInfoFromRedis(groupId string) *GroupInfo {

	var groupInfo GroupInfo
	_ = GetJSON(groupRelayGroupIdRedisKey+groupId, &groupInfo)

	return &groupInfo
}

func getGroupRelayMembersFromRedis(groupId string) ([]string, error) {
	result, err := GetHashList(groupRelayMembersRedisKey, groupId)

	if err != nil {
		return nil, err
	}

	return result, nil
}
