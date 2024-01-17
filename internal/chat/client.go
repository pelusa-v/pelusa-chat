package chat

type Message struct {
	IdOrigin      string
	IdDestination string
}

type Client struct {
	Id int
}

func (c *Client) SendMessage() {

}
