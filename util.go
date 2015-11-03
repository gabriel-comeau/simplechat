package main

import (
	"fmt"
	"strings"
	"sync"
)

var (
	nextId int64
	lock   *sync.Mutex
)

func init() {
	lock = new(sync.Mutex)
	nextId = 1
}

func getNextId() int64 {
	lock.Lock()
	defer lock.Unlock()
	retVal := nextId
	nextId++
	return retVal
}

// Formats a message's test in the "nickname: message" format
func formatBroadCastMessage(message string, sender *client) string {
	message = strings.Trim(message, " \n")
	output := ""
	if sender.username != "" {
		output = fmt.Sprintf("%v: %v\n", sender.username, message)
	} else {
		output = fmt.Sprintf("%v: %v\n", sender.id, message)
	}

	return output
}

// Formats a message in format suitable for private messages
func formatWhisperMessage(message string, sender *client) string {
	message = strings.Trim(message, " \n")
	output := ""
	if sender.username != "" {
		output = fmt.Sprintf("<PRIVATE MESSSAGE> %v: %v\n", sender.username, message)
	} else {
		output = fmt.Sprintf("<PRIVATE MESSSAGE> %v: %v\n", sender.id, message)
	}

	return output
}

func formatErrorMessage(message string) string {
	return fmt.Sprintf("<ERROR>: %v\n", message)
}
