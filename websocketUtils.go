package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

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

func (c *WebSocketClient) SendMessage(message string) error {
	return c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *WebSocketClient) SendPing() error {
	return c.Conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *WebSocketClient) ReceiveMessage() {
	defer c.Conn.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
		go c.HandlerMessage(message)
	}
}

func SendHeartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		go SendPingAndReconnect()
	}
}

func SendPingAndReconnect() {
	mu.RLock() //add read lock
	defer mu.RUnlock()

	if len(clientsMap) > 0 {
		for _, client := range clientsMap {
			err := client.SendMessage("ping")
			if err != nil {
				log.Printf("Failed to send heartbeat to %s: %v", client.MonitorgRelaysInfo.RelayUrl, err)
				mu.RUnlock()
				ConnectToInitRelays([]MonitoringRelaysInfo{*client.MonitorgRelaysInfo})
				mu.RLock()
			}
		}
	}
}

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

		err = c.SendMessage(authMsg)
		if err != nil {
			log.Printf("send auth msg to relay:%s failed.authMsg=:%s", c.MonitorgRelaysInfo.RelayUrl, authMsg)
		}
	}
}

func HandleInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Println("Interrupt signal received, closing connections...")
	for _, value := range clientsMap {
		value.Close()
	}
	os.Exit(0)
}

// ConnectToInitRelays connects to all relays and initializes clients.
func ConnectToInitRelays(relays []MonitoringRelaysInfo) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for _, relay := range relays {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(relay MonitoringRelaysInfo) {
			defer wg.Done()
			defer func() { <-semaphore }()

			client, err := NewWebSocketClient(&relay)
			if err != nil {
				log.Printf("Failed to connect to WebSocket server %s: %v", relay.RelayUrl, err)
				return
			}

			log.Printf("Connected to %s successfully!", relay.RelayUrl)

			clients = append(clients, client)
			clientsMap[relay.RelayUrl] = client

			go client.ReceiveMessage()

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
		}(relay)
	}

	wg.Wait()
	log.Println("All relays have been initialized and clients are connected.")
}
