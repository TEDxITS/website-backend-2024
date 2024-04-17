package websocket

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type (
	TicketQueue interface {
		Serve(*gin.Context)
		WaitQueueTurn(*Client) error

		Register(*Client)
		Unregister(*Client)
		SendOperation(operation)
	}

	ticketQueue struct {
		Hub        QueueHub
		jwtService service.JWTService
	}
)

func NewTicketQueue(hub QueueHub, jwt service.JWTService) TicketQueue {
	return &ticketQueue{
		Hub:        hub,
		jwtService: jwt,
	}
}

func (Handle *ticketQueue) Serve(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if nil != err {
		return
	}
	defer conn.Close()

	client := NewClient(conn)

	for !client.IsAuthenticated() {
		var token string

		message, err := client.ReadMessage()
		if nil != err {
			return
		}

		if !strings.HasPrefix(string(message), constants.CTX_KEY_TOKEN) {
			if err := client.SendTextMessage(dto.ErrWSInvalidCommand.Error()); nil != err {
				return
			}
			continue
		}

		if _, err := fmt.Sscanf(string(message), dto.WSOCKET_AUTH_REQUEST, &token); nil != err {
			if err := client.SendTextMessage(err.Error()); nil != err {
				return
			}
			continue
		}

		userId, _, err := Handle.jwtService.GetPayloadInsideToken(token)
		if nil != err {
			if err := client.SendTextMessage(dto.ErrWSInvalidToken.Error()); nil != err {
				return
			}
			continue
		}

		client.SetUserID(userId)
		client.SetAuthenticated()

		if err := client.SendTextMessage(dto.WSOCKET_AUTH_SUCCESS); nil != err {
			return
		}
	}

	// already in queue, (detect using one account to queue multiple times)
	if Handle.Hub.IsInQueueByUserID(client.UserID) {
		_ = client.SendTextMessage(dto.ErrWSAlreadyInQueue.Error())
		return
	}

	// register client to hub for tracking
	Handle.Register(client)
	defer Handle.Unregister(client)

	time.Sleep(200 * time.Millisecond)
	if err := Handle.WaitQueueTurn(client); nil != err {
		return
	}

	if err := client.SendTextMessage(dto.WSOCKET_TRANSACTION_START); nil != err {
		return
	}

	incoming := make(chan []byte)
	go func() {
		for {
			// otherwise, mutex block, thread hang
			_, message, err := client.Conn.ReadMessage()
			if nil != err {
				client.Done(err)
				return
			}
			incoming <- message
		}
	}()

	for {
		select {
		case message := <-incoming:
			if len(message) == 0 {
				if err := client.SendTextMessage(dto.ErrWSInvalidCommand.Error()); nil != err {
					return
				}
				continue
			}

			switch string(message) {
			case dto.WSOCKET_ENUM_WITH_MERCH_REQUEST:
				Handle.SendOperation(operation{
					client:  client,
					command: dto.WSOCKET_ENUM_WITH_MERCH_REQUEST,
				})
			case dto.WSOCKET_ENUM_NO_MERCH_REQUEST:
				Handle.SendOperation(operation{
					client:  client,
					command: dto.WSOCKET_ENUM_NO_MERCH_REQUEST,
				})
			default:
				if err := client.SendTextMessage(dto.ErrWSInvalidCommand.Error()); nil != err {
					return
				}
			}
		case message := <-client.Notification:
			if err := client.SendTextMessage(message); err != nil {
				return
			}
		case err := <-client.Quit:
			if err == nil {
				client.SendTextMessage(dto.WSOCKET_TRANSACTION_SUCCESS)
			}

			return
		}
	}

}

func (Handle *ticketQueue) WaitQueueTurn(client *Client) error {
	if !client.IsWaiting() {

		// if not buffered at construct, we need to consume or else
		// the thread will hang/block
		// <-client.Notification

		return nil
	}

	queueNumber := Handle.Hub.GetWaitingLength()
	if err := client.SendTextMessage(fmt.Sprintf(dto.WSOCKET_QUEUE_NUMBER, queueNumber)); nil != err {
		return err
	}

	for {
		select {
		case next := <-client.Next:
			if next == client {
				return nil
			}

			queueNumber--
			if err := client.SendTextMessage(fmt.Sprintf(dto.WSOCKET_QUEUE_NUMBER, queueNumber)); nil != err {
				return err
			}
		case notif := <-client.Notification:
			if notif == dto.ErrWSMainEventFull.Error() {
				client.SendTextMessage(dto.ErrWSMainEventFull.Error())
				return dto.ErrWSMainEventFull
			}
		}
	}
}

func (Handle *ticketQueue) Register(client *Client) {
	Handle.Hub.GetRegisterChannel() <- client
}

func (Handle *ticketQueue) Unregister(client *Client) {
	Handle.Hub.GetUnregisterChannel() <- client
}

func (Handle *ticketQueue) SendOperation(ops operation) {
	Handle.Hub.GetOperationChannel() <- ops
}
