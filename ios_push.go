package main

import (
	"fmt"
	"strings"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

func InitClientMap() {
	log.Println("Begin Init APNS Client!!!")
	bundleIds := config.IOSPushInfo.BundleIds
	p12paths := config.IOSPushInfo.P12Pathes
	passwords := config.IOSPushInfo.Passwords

	bundleArr := strings.Split(bundleIds, ",")
	p12Arr := strings.Split(p12paths, ",")
	passwordArr := strings.Split(passwords, ",")

	for i := 0; i < len(bundleArr); i++ {
		AddClient(bundleArr[i], p12Arr[i], passwordArr[i])
	}
	log.Println("Init APNS Client Successful!!!")
}

// AddClient adds a new APNs client to the map
func AddClient(name string, certFile string, password string) error {
	cert, err := certificate.FromP12File(certFile, password)
	if err != nil {
		return err
	}
	client := apns2.NewClient(cert).Production()
	clientMap[name] = client
	return nil
}

type IosPushInitDTO struct {
	DeviceToken string `json:"deviceToken"`
	BundleId    string `json:"bundleid"`
}

func GetDeviceInfoByToken(deviceToken string) IosPushInitDTO {
	var iosPushInitDTO IosPushInitDTO
	for key := range clientMap {
		if strings.HasPrefix(deviceToken, key) {
			iosPushInitDTO.DeviceToken = deviceToken[len(key):]
			iosPushInitDTO.BundleId = key
		}
	}

	return iosPushInitDTO
}

// Push sends a notification to the specified device tokens using the specified client
func PushIos(deviceToken string, title, message string) error {
	iosPushInitDTO := GetDeviceInfoByToken(deviceToken)
	client, exists := clientMap[iosPushInitDTO.BundleId]
	if !exists {
		return fmt.Errorf("client %s not found", iosPushInitDTO.BundleId)
	}

	notification := &apns2.Notification{}
	notification.Topic = iosPushInitDTO.BundleId
	notification.Payload = payload.NewPayload().AlertTitle(title).AlertBody(message).Badge(1)
	notification.DeviceToken = iosPushInitDTO.DeviceToken

	res, err := client.Push(notification)
	if err != nil {
		log.Printf("Error sending to %s: %v", iosPushInitDTO.DeviceToken, err)
	} else {
		log.Printf("Sent to %s: %v %v\n", iosPushInitDTO.DeviceToken, res.StatusCode, res.ApnsID)
	}

	return nil
}
