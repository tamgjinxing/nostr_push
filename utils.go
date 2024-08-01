package main

import "encoding/json"

// contains 检查切片是否包含某个元素
func Contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// anyMatch 检查是否有任何一个元素在目标切片中
func AnyMatch(targets, kinds []int) bool {
	for _, target := range targets {
		if Contains(kinds, target) {
			return true
		}
	}
	return false
}

// RemoveDuplicates 去重函数
func RemoveDuplicates(strSlice []string) []string {
	// 创建一个 map 来跟踪唯一的字符串
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string

	// 遍历原始字符串切片
	for _, item := range strSlice {
		// 如果字符串尚未在 map 中存在
		if _, exists := uniqueMap[item]; !exists {
			// 将字符串添加到 map 中
			uniqueMap[item] = true
			// 也将字符串添加到唯一的字符串切片中
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice
}

func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
