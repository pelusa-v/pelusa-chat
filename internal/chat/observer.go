package chat

type Observer struct {
	Clients               []*Client
	SubscribeClientChan   chan *Client
	UnsubscribeClientChan chan *Client
}

func NewObserver() *Observer {
	return &Observer{
		Clients:               make([]*Client, 0),
		SubscribeClientChan:   make(chan *Client),
		UnsubscribeClientChan: make(chan *Client),
	}
}

func (o *Observer) Start() {
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
		}
	}
}
