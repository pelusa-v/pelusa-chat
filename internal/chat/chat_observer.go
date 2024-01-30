package chat

type ChatObserver struct {
	Clients               []*Client
	SubscribeClientChan   chan *Client
	UnsubscribeClientChan chan *Client
	Broadcast             chan *string
}

func NewChatObserver() *ChatObserver {
	return &ChatObserver{
		Clients:               make([]*Client, 0),
		SubscribeClientChan:   make(chan *Client),
		UnsubscribeClientChan: make(chan *Client),
		Broadcast:             make(chan *string),
	}
}

func (o *ChatObserver) Start() {
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
			// case channel := <-o.Broadcast:

		}
	}
}
