package chat

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
}

type ClientJson struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Observer *Observer
}

func NewClient(id string, name string, obs *Observer) *Client {
	return &Client{
		Id:       id,
		Name:     name,
		Observer: obs,
	}
}

func (c *Client) SendMessage() {

}
