package main

import (
	"strings"
)

type NotificationBean struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

type MessageBean struct {
	Data         map[string]string `json:"data"`
	Notification *NotificationBean `json:"notification"`
}

func NewMessageBean() *MessageBean {
	return &MessageBean{
		Data: make(map[string]string),
	}
}

func (m *MessageBean) PutData(key string, value string) *MessageBean {
	m.Data[key] = value
	return m
}

func (m *MessageBean) SetData(data map[string]string) *MessageBean {
	m.Data = data
	return m
}

func (m *MessageBean) SetNotification(notification *NotificationBean) *MessageBean {
	m.Notification = notification
	return m
}

func sendPush(deviceToken string, messageBean *MessageBean) {
	if deviceToken == "" {
		return
	}

	log.Printf("before to start send!, url: %s, message: %s", deviceToken, ToJSON(messageBean))

	callUrl := deviceToken

	if strings.Contains(strings.ToLower(deviceToken), "embedded-fcm/fcm") {
		callUrl = "https://www.0xchat.com" + deviceToken[strings.Index(deviceToken, "/FCM"):]
	}

	if strings.Contains(strings.ToLower(deviceToken), "ntfy") || strings.Contains(strings.ToLower(deviceToken), "conversations") {
		callUrl = deviceToken
	}

	response, statusCode, err := HttpClient.Post(callUrl, nil, messageBean)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Printf("Status Code: %d, Response Body: %s\n", statusCode, response)
}

func PushAndroid(deviceId, title, message string, isCallPush bool, groupId string) {
	var messageBean *MessageBean
	if groupId == "" {
		if isCallPush {
			messageBean = NewMessageBean().
				PutData("msgType", "1").
				SetNotification(&NotificationBean{
					Title: title,
					Body:  message,
				})
		} else {
			messageBean = NewMessageBean().
				SetNotification(&NotificationBean{
					Title: title,
					Body:  message,
				})
		}
	} else {
		if isCallPush {
			messageBean = NewMessageBean().
				PutData("msgType", "1").
				PutData("groupId", groupId).
				SetNotification(&NotificationBean{
					Title: title,
					Body:  message,
				})
		} else {
			messageBean = NewMessageBean().
				PutData("groupId", groupId).
				SetNotification(&NotificationBean{
					Title: title,
					Body:  message,
				})
		}
	}

	sendPush(deviceId, messageBean)
}
