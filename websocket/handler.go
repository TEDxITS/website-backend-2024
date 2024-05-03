package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
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
		jwtService config.JWTService
	}
)

func NewTicketQueue(hub QueueHub, jwt config.JWTService) TicketQueue {
	return &ticketQueue{
		Hub:        hub,
		jwtService: jwt,
	}
}

func (Handler *ticketQueue) Serve(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if nil != err {
		return
	}
	defer conn.Close()

	client := NewClient(conn)

	// 10s to authenticate
	client.SetDeadline(time.Now().Add(constants.WSOCKET_AUTH_TIME_LIMIT))
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

		userId, _, err := Handler.jwtService.GetPayloadInsideToken(token)
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
	if Handler.Hub.IsInQueueByUserID(client.UserID) {
		_ = client.SendTextMessage(dto.ErrWSAlreadyInQueue.Error())
		return
	}

	// register client to hub for tracking
	Handler.Register(client)

	// set no timeout when waiting in queue
	client.SetDeadline(time.Time{})

	// the only error comes from this queue turn
	// is only if the main event is full
	time.Sleep(200 * time.Millisecond)
	if err := Handler.WaitQueueTurn(client); nil != err {
		if err != dto.ErrWSMainEventFull {
			Handler.Unregister(client)
		}
		client.Conn.Close()
		return
	}

	// we don't wanna unregister if the the event is full
	// because it would trigger Hub.BroadcastNext()
	// which in turn will trigger client.Notify()
	// which is a waste of resource

	// the two mentioned method above is already executed
	// in the event of the main event is full

	// so we only unregister in case of:
	// 	1. client's transaction is  done
	//  2. client's transaction timeout
	defer Handler.Unregister(client)

	// 3m to finish transaction
	client.SetDeadline(time.Now().Add(constants.WSOCKET_TRANSACTION_TIME_LIMIT))
	if err := client.SendTextMessage(dto.WSOCKET_TRANSACTION_START); nil != err {
		return
	}

	withMerchPrice, noMerchPrice := Handler.Hub.GetBasePrice()

	message, _ := json.Marshal(struct {
		NoMerchBasePrice   int `json:"no_merch_price"`
		WithMerchBasePrice int `json:"with_merch_price"`
	}{
		withMerchPrice,
		noMerchPrice,
	})
	if err := client.SendTextMessage(string(message)); nil != err {
		return
	}

	messageFromClient := make(chan []byte)
	go func() {
		for {
			// otherwise, mutex block, thread hang
			_, message, err := client.Conn.ReadMessage()
			if nil != err {
				client.Done(err)
				return
			}
			messageFromClient <- message
		}
	}()

	for {
		select {
		// case for the client's action:
		// 1. change the transaction type (with/without merch)
		// 2. seat selection
		case message := <-messageFromClient:
			if len(message) == 0 {
				if err := client.SendTextMessage(dto.ErrWSInvalidCommand.Error()); nil != err {
					return
				}
				continue
			}

			switch string(message) {
			case dto.WSOCKET_ENUM_WITH_MERCH_REQUEST:
				Handler.SendOperation(operation{
					client:  client,
					command: dto.WSOCKET_ENUM_WITH_MERCH_REQUEST,
				})
			case dto.WSOCKET_ENUM_NO_MERCH_REQUEST:
				Handler.SendOperation(operation{
					client:  client,
					command: dto.WSOCKET_ENUM_NO_MERCH_REQUEST,
				})
			default:
				if err := client.SendTextMessage(dto.ErrWSInvalidCommand.Error()); nil != err {
					return
				}
			}

		// notify client of changing in information such as the stock of merch or the seat
		case messageFromHub := <-client.Notification:
			if err := client.SendTextMessage(messageFromHub); err != nil {
				return
			}

		// notify the handler to exit/return in case of the transaction is done
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

	// get and send the initial queue line number
	queueNumber := Handle.Hub.GetWaitingLength()
	message, _ := json.Marshal(dto.S2CQueueLineInfo{
		QueueNumber: queueNumber,
	})
	if err := client.SendTextMessage(string(message)); nil != err {
		return err
	}

	for {
		select {
		// notification of the main event is full
		case notif := <-client.Notification:
			if notif == dto.ErrWSMainEventFull.Error() {
				client.SendTextMessage(notif)
				return dto.ErrWSMainEventFull
			}
		// receiving the next client to forward to transaction
		// if not the current client, then it will simply decrement its queue number
		case next := <-client.Next:
			if next == client {
				return nil
			}

			queueNumber--
			message, _ := json.Marshal(dto.S2CQueueLineInfo{
				QueueNumber: queueNumber,
			})

			if err := client.SendTextMessage(string(message)); nil != err {
				return err
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
