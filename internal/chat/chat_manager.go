package chat

import "fmt"

type ChatManager struct {
	Clients                   []*Client
	SubscribeClientChan       chan *Client
	UnsubscribeClientChan     chan *Client
	BroadcastNotificationChan chan *WebSocketData
	SendMessageChan           chan *WebSocketData
}

func (manager *ChatManager) Start() {
	for {
		select {
		case channel := <-manager.SubscribeClientChan:
			manager.Clients = append(manager.Clients, channel)

		case channel := <-manager.UnsubscribeClientChan:
			for i, client := range manager.Clients {
				if client.Id == channel.Id {
					manager.Clients = append(manager.Clients[:i], manager.Clients[i+1:]...)
				}
			}

		case channel := <-manager.SendMessageChan: // send message to destination client
			for _, client := range manager.Clients {
				if client.Id == channel.Message.IdDestination {
					fmt.Println("Sending message to " + channel.Message.IdDestination)
					fmt.Println("Content " + channel.Message.Content)
					client.ReceiveMessageChan <- channel
				}
			}

		case channel := <-manager.BroadcastNotificationChan: // send notification to destination client
			for _, client := range manager.Clients {
				fmt.Println("Sending message to " + channel.Notification.ClientName)
				fmt.Println("Content " + channel.Notification.Content)
				client.ReceiveMessageChan <- channel
			}
		}
	}
}

var Manager = ChatManager{
	Clients:                   make([]*Client, 0),
	SubscribeClientChan:       make(chan *Client),
	UnsubscribeClientChan:     make(chan *Client),
	BroadcastNotificationChan: make(chan *WebSocketData),
	SendMessageChan:           make(chan *WebSocketData),
}
