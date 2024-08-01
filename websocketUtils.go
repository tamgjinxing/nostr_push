package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

type HanlderEventInfo struct {
	Event  *nostr.Event
	Client *WebSocketClient
}

type SubKindsLimit struct {
	SubscribeKinds []int
	Limit          int
}

type MonitoringRelaysInfo struct {
	RelayUrl       string
	SubKindsLimits []SubKindsLimit
	GroupRelayFlag bool
}

// WebSocketClient 结构体包含 WebSocket 连接
type WebSocketClient struct {
	Conn               *websocket.Conn
	Challenge          string
	AuthedPublicKey    string
	MonitorgRelaysInfo *MonitoringRelaysInfo
}

// NewWebSocketClient 创建新的 WebSocketClient 实例并连接到 WebSocket 服务
func NewWebSocketClient(monitorRelayInfo *MonitoringRelaysInfo) (*WebSocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(monitorRelayInfo.RelayUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &WebSocketClient{
		Conn:               conn,
		MonitorgRelaysInfo: monitorRelayInfo,
	}
	return client, nil
}

// SendMessage 发送消息到 WebSocket 连接
func (c *WebSocketClient) SendMessage(message string) error {
	return c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// ReceiveMessage 接收来自 WebSocket 连接的消息
func (c *WebSocketClient) ReceiveMessage() {
	defer c.Conn.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		c.HandlerMessage(message)
	}
}

// SendHeartbeat 定期发送心跳消息
func (c *WebSocketClient) SendHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		err := c.SendMessage("ping")
		if err != nil {
			log.Printf("Failed to send heartbeat to %s: %v", c.MonitorgRelaysInfo.RelayUrl, err)
			return
		}
	}
}

// Close 关闭 WebSocket 连接
func (c *WebSocketClient) Close() {
	c.Conn.Close()
}

func (c *WebSocketClient) HandlerMessage(message []byte) {
	log.Printf("Received message from %s: message:%s\n", c.MonitorgRelaysInfo.RelayUrl, message)

	var rawMessage []json.RawMessage
	err := json.Unmarshal(message, &rawMessage)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
	}

	// Unmarshal the type
	var eventType string
	err = json.Unmarshal(rawMessage[0], &eventType)
	if err != nil {
		log.Printf("Error unmarshaling event type: %v", err)
	}

	if eventType == "EVENT" {
		var payload nostr.Event
		err = json.Unmarshal(rawMessage[2], &payload)
		if err != nil {
			log.Printf("Error unmarshaling payload: %v", err)
		}

		hanlderEventInfo := HanlderEventInfo{
			Event:  &payload,
			Client: c,
		}

		if Contains(GroupRelayChannelKinds, payload.Kind) {
			groupRelayChannel <- hanlderEventInfo
		} else if Contains(MonmentChannelKinds, payload.Kind) {
			monentChannel <- hanlderEventInfo
		} else if payload.Kind == 9000 {
			inviteToGroupChannel <- hanlderEventInfo
		} else {
			eventChannel <- hanlderEventInfo
		}
	}

	if eventType == "AUTH" {
		var challenge string
		err = json.Unmarshal(rawMessage[1], &challenge)
		if err != nil {
			log.Printf("Error unmarshaling event type: %v", err)
		}
		c.AuthedPublicKey = config.PushBotInfo.PublicKey
		c.Challenge = challenge

		authMsg, err := GenerateAuthMsg(c.MonitorgRelaysInfo.RelayUrl, challenge)
		if err != nil {
			log.Printf("generateAuthMsg failed:%v\n", err)
		}

		log.Printf("authMsg:%s\n", authMsg)
		err = c.SendMessage(authMsg)
		if err != nil {
			log.Printf("send auth msg to relay:%s failed.authMsg=:%s", c.MonitorgRelaysInfo.RelayUrl, authMsg)
		}
	}
}

// HandleInterrupt 捕获中断信号，优雅地关闭 WebSocket 连接
func HandleInterrupt(clients []*WebSocketClient) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Println("Interrupt signal received, closing connections...")
	for _, client := range clients {
		client.Close()
	}
	os.Exit(0)
}

func ConnectToInitRelays(relays []MonitoringRelaysInfo) {
	// 连接到所有 WebSocket 服务
	for _, relay := range relays {
		client, err := NewWebSocketClient(&relay)
		if err != nil {
			log.Printf("Failed to connect to WebSocket server %s: %v", relay.RelayUrl, err)
			continue
		}

		client.MonitorgRelaysInfo = &relay
		clients = append(clients, client)
	}

	// 启动接收消息的 goroutine
	for _, client := range clients {
		go client.ReceiveMessage()
		go client.SendHeartbeat(10 * time.Second) // 每10秒发送一次心跳

		go func(c *WebSocketClient) {
			SubKindsLimits := c.MonitorgRelaysInfo.SubKindsLimits
			if len(SubKindsLimits) > 0 {
				for _, subKindsLimit := range SubKindsLimits {
					subId := GenerateRandomString(32)
					filters := map[string]interface{}{
						"kinds": subKindsLimit.SubscribeKinds,
						"limit": subKindsLimit.Limit,
					}
					reqMsg := GenerateSubscribeMsg(subId, filters)
					err := c.SendMessage(reqMsg)
					if err != nil {
						log.Printf("Failed to send message to %s: %v", c.MonitorgRelaysInfo.RelayUrl, err)
					} else {
						log.Printf("Sent message to %s: %s\n", c.MonitorgRelaysInfo.RelayUrl, reqMsg)
					}
				}
			}

			select {}
		}(client)
	}
}
