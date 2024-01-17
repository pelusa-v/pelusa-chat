package main

import "github.com/pelusa-v/pelusa-chat.git/internal/chat"

func main() {
	obs := chat.NewObserver()
	go obs.Start()
}
