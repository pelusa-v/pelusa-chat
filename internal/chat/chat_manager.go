package chat

import "fmt"

type ChatManager struct {
	Clients               []*Client
	SubscribeClientChan   chan *Client
	UnsubscribeClientChan chan *Client
	BroadcastMessageChan  chan *string
	SendMessageChan       chan *Message
}

func (o *ChatManager) Start() {
	for {
		select {
		case channel := <-o.SubscribeClientChan:
			o.Clients = append(o.Clients, channel)

		case channel := <-o.UnsubscribeClientChan:
			for i, client := range o.Clients {
				if client.Id == channel.Id {
					o.Clients = append(o.Clients[:i], o.Clients[i+1:]...)
				}
			}

		case channel := <-o.SendMessageChan: // send message to destination client
			for _, client := range o.Clients {
				if client.Id == channel.IdDestination {
					fmt.Println("Sending message to " + channel.IdDestination)
					fmt.Println("Content " + channel.Content)
					client.ReceiveMessageChan <- channel.Content
				}
			}
		}
	}
}

var Manager = ChatManager{
	Clients:               make([]*Client, 0),
	SubscribeClientChan:   make(chan *Client),
	UnsubscribeClientChan: make(chan *Client),
	BroadcastMessageChan:  make(chan *string),
	SendMessageChan:       make(chan *Message),
}
