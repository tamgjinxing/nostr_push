package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// 生成指定长度的随机字符串
func GenerateRandomString1(length int) string {
	// 定义字符集合
	charset := "ABCDEFabcdef0123456789"

	// 生成随机字符串
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}

	return string(randomString)
}

// 生成relay的订阅信息（req）
// reqId reqId
func GenerateSubscribeMsg1(reqId string, filters map[string]interface{}) string {
	jsonMap := make(map[string]interface{})
	for key, value := range filters {
		jsonMap[key] = value
	}

	jsonArray := []interface{}{"REQ", reqId, jsonMap}

	jsonString, err := json.Marshal(jsonArray)
	if err != nil {
		return "" // 返回空字符串或者其他错误处理逻辑
	}

	return string(jsonString)
}

// 生成退订消息
func GenerateUnSubscribeMsg1(reqId string) string {
	jsonArray := []interface{}{"CLOSE", reqId}

	jsonString, err := json.Marshal(jsonArray)
	if err != nil {
		return "" // 返回空字符串或者其他错误处理逻辑
	}

	return string(jsonString)
}

func Gen5xxxEvent(content string, tags []Tag) (string, error) {
	event := Event{
		CreatedAt: time.Now().Unix(),
		Kind:      24133,
		Content:   "Ap/ZsXS5w7DnLDYxrcickhkKlE1kdKSWDzLzxpwW+Pe100Q+HId/+pvOVgGQJlQbrIu/a869Vmv9xCX2O+rRkfZ4czacAAZDhDN64gNDHaUrpQY5AmNaJiI95QMVzGwIDdofgZVcrtYOfTm91jDUdKyrK5OczF62qV77oARtzEdUJuRk8bbrMDLwvE7CPWGis3Uxre1be1KU+cKp3fZ/Mf8AhI2ETpBIPzxh9zHo6h73R13CjDjUvRHPcmtTMaxQv9VT4lkUC1g/wwKXJDtUs1geiTjaIzpC3hqmoxaUC9hBgZNbjJVpHvtIGsMMhOk8p6R/jWbMdEsdz2FZOjLMZe+OiRIhnumUrHCqIGVxThm2AS1irmY6whG5JbFe2To34o75",
		Tags:      tags,
	}
	event.Sign("52bb1a30d3ee0dc7b82b6b9f88fff743fdac4d1b8baecdeca41fb1ede50e83ed")
	return event.ToEventString()
}

func Gen9735Event(content string, tags []Tag) (string, error) {
	event := Event{
		CreatedAt: 1734680216,
		Kind:      9735,
		Content:   "",
		Tags:      tags,
	}

	//pubkey: 3964b75fe5d92e71d87b93202a30a316bcd957c55157b3ddea7d5c4a8c7d5e45
	//prikey: 659ebf89f4bcef6536a5f3bbb3a55bcb7c579e6325ccf1d64ced569cf476d756
	event.Sign("659ebf89f4bcef6536a5f3bbb3a55bcb7c579e6325ccf1d64ced569cf476d756")
	return event.ToEventString()
}

func TestGen9735Event() {
	tags := make([]Tag, 0, 4) // 初始长度为0，容量为5

	p := []string{"p", "093dff31a87bbf838c54fd39ff755e72b38bd6b7975c670c0f2633fa7c54ddd0"}
	P := []string{"P", "aaf092aabda3304cf5f5a5c8c0717795a34a7ab745a124eae07fc4820a8f04a3"}
	bolt11 := []string{"bolt11", "lnbc40n1pnk8kp3dqjgfjhxapqwa5hx6r9wvnp4q05dsfskvh5j70j8v2fkyvc88w9u7fxfnrr2zv2x8y7mzwuzcmg3gpp55jl48lay9748vht5knaa0xvphp3m7a2vk6m3eracdht2rjy93d3qsp5uvkftjk6dmr8p38rvwlacvmlms82wdpzt2c8kfz38e8j8cjlmvyq9qyysgqcqpcxqyz5vqrzjqw9fu4j39mycmg440ztkraa03u5qhtuc5zfgydsv6ml38qd4azymlapyqqqqqqqxucqqqqlgqqqq86qqjqzqzjt8efzcdg7gswqchy7fwsvcmyj5a80lt3k39htpw476vrh8x3zfwm6wavgh929c0xtaa32u5kxa3fpcfh3y6a6h79yzeyuc692lqql67znk"}
	preimage := []string{"preimage", "37f94fdd60e89eca5b38cda2a5b256faefddf8dd7dcd7aa8518a92ee312e103e"}
	description := []string{"description", "{\"id\":\"4dc8d9b0aa635665992076b0b190750475b90e0c2b94b16b2d3bb0e10254694a\",\"pubkey\":\"aaf092aabda3304cf5f5a5c8c0717795a34a7ab745a124eae07fc4820a8f04a3\",\"created_at\":1734596656,\"kind\":9734,\"tags\":[[\"relays\",\"wss://relay.0xchat.com\"],[\"amount\",\"4000\"],[\"lnurl\",\"lnurl1dp68gurn8ghj7em9w3skccne9e3k7mf09emk2mrv944kummhdchkcmn4wfk8qtmpd3kx2m3sx5cnz90uh99\"],[\"p\",\"392d3fcca8e3e924e256625983268608926c9a54e0703c9b3fffae320ad1c86a\"],[\"e\",\"c2886303e64a6fa97ab555574e2aa5392f774aa2e13202ee32f1024d332dfc4c\"]],\"content\":\"Best wishes\",\"sig\":\"25444e30182636de3e23e33294f0a070234e0c100794e746c4835dafbf016bc6354f234cab496c959e5358f5c50042641d3309f0d92bab52db099e5c5e8ba93c\"}"}

	tags = append(tags, p)
	tags = append(tags, P)
	tags = append(tags, bolt11)
	tags = append(tags, preimage)
	tags = append(tags, description)

	eventString, err := Gen9735Event("", tags)
	if err != nil {
		return
	}

	fmt.Println(eventString)
}

func main1() {
	TestGen9735Event()
}
