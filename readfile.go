package main

import (
	"bufio"
	"os"
)

func ReadRelaysFile(relaysFileName string) {
	log.Printf("relaysFileName:%s\n", relaysFileName)
	file, err := os.Open(relaysFileName)
	if err != nil {
		log.Printf("can not open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read file failed: %v", err)
	}

	needMonitorRelays = append(needMonitorRelays, lines...)

}

func GetMonitorChatRelays(topN int) []string {
	var firstNLines []string
	if len(needMonitorRelays) > topN {
		firstNLines = needMonitorRelays[:topN]
	} else {
		firstNLines = needMonitorRelays
	}

	return firstNLines

}
