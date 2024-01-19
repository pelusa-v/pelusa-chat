package chat

import "github.com/gofiber/contrib/websocket"

type Message struct {
	IdOrigin      string `json:"id_origin"`
	IdDestination string `json:"id_destination"`
}

type ClientJson struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Observer      *Observer
	WebsocketConn *websocket.Conn
}

func NewClient(id string, name string, obs *Observer, conn *websocket.Conn) *Client {
	return &Client{
		Id:            id,
		Name:          name,
		Observer:      obs,
		WebsocketConn: conn,
	}
}

func (c *Client) SendMessage() {

}
