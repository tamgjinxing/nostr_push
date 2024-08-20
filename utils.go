package main

import "encoding/json"

func Contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func AnyMatch(targets, kinds []int) bool {
	for _, target := range targets {
		if Contains(kinds, target) {
			return true
		}
	}
	return false
}

func RemoveDuplicates(strSlice []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string

	for _, item := range strSlice {
		if _, exists := uniqueMap[item]; !exists {
			uniqueMap[item] = true
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
