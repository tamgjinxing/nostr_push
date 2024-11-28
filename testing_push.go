package main

func main11() {
	// Initialize push management
	pushManager = NewPushManager(androidPusher, iosPusher)

	configName := "config.json"

	// Reading configuration files
	err := ReadConfig(configName)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	// For ios push
	InitClientMap()

	deviceToken := "com.0xchat.app29843a7e252caf1526ae8224ac384f283f5bc50d688c983f38f5aabf55d555ea"
	title := "testpush-20241128-001-title"
	message := "testpush-20241128-001-message"
	isCallPush := false
	groupId := ""

	log.Println("Testing ios push with calling mp3")

	PushIos(deviceToken, title, message, isCallPush, groupId)
}
