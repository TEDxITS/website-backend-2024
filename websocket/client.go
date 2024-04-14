package websocket

import (
	"sync"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/gorilla/websocket"
)

type (
	Client struct {
		Conn *websocket.Conn

		UserID        string
		Authenticated bool
		Waiting       bool
		WithMerch     bool

		Next         chan *Client
		Notification chan string
		Quit         chan error

		*sync.Mutex
	}
)

func NewClient(conn *websocket.Conn) *Client {
	deadline := time.Now().Add(constants.WSOCKET_TIME_LIMIT)
	conn.SetReadDeadline(deadline)
	conn.SetWriteDeadline(deadline)

	return &Client{
		Conn: conn,

		Authenticated: false,
		Waiting:       false,
		WithMerch:     false,

		Next:         make(chan *Client, 1),
		Notification: make(chan string),
		Quit:         make(chan error),

		Mutex: new(sync.Mutex),
	}
}

func (c *Client) SendTextMessage(msg string) error {
	c.Lock()
	defer c.Unlock()

	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); nil != err {
		return err
	}
	return nil
}

func (c *Client) ReadMessage() ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	_, message, err := c.Conn.ReadMessage()
	if nil != err {
		return nil, err
	}

	return message, nil
}

func (c *Client) Done(err error) {
	c.Quit <- err
}

func (c *Client) Notify(notif string) {
	c.Notification <- notif
}

func (c *Client) InformNext(next *Client) {
	c.Next <- next
}

func (c *Client) SetUserID(userID string) {
	c.Lock()
	defer c.Unlock()

	c.UserID = userID
}

func (c *Client) SetAuthenticated() {
	c.Lock()
	defer c.Unlock()

	c.Authenticated = true
}

func (c *Client) IsAuthenticated() bool {
	return c.Authenticated
}

func (c *Client) SetWaiting() {
	c.Lock()
	defer c.Unlock()

	c.Waiting = true
}

func (c *Client) IsWaiting() bool {
	return c.Waiting
}

func (c *Client) SetWithMerch(b bool) {
	c.Lock()
	defer c.Unlock()

	c.WithMerch = b
}

func (c *Client) IsWithMerch() bool {
	return c.WithMerch
}
