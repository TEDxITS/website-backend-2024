package websocket

import (
	"encoding/json"
	"sync"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/repository"
)

type (
	QueueHub interface {
		PushWaiting(*Client)
		WalkRemove(*Client)
		UpdateMaxTransaction()
		ClientDone(string)

		BroadcastNext()
		BroadcastStock()

		IsInQueueByUserID(string) bool
		GetWaitingLength() int

		GetRegisterChannel() chan *Client
		GetUnregisterChannel() chan *Client
		GetOperationChannel() chan operation
	}

	queueHub struct {
		Transaction    []*Client
		Waiting        []*Client
		MaxTransaction int

		Register   chan *Client
		Unregister chan *Client
		Operation  chan operation
		Done       chan string

		NoMerchID   string
		WithMerchID string

		repository repository.EventRepository

		*sync.Mutex
	}

	operation struct {
		client  *Client
		command string
	}
)

func RunConnHub(repo repository.EventRepository, max int, noMerchID, withMerchID string) QueueHub {
	hub := &queueHub{
		Transaction:    make([]*Client, 0),
		Waiting:        make([]*Client, 0),
		MaxTransaction: max,

		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Operation:  make(chan operation),
		Done:       make(chan string),

		NoMerchID:   noMerchID,
		WithMerchID: withMerchID,

		repository: repo,

		Mutex: new(sync.Mutex),
	}

	go func() {
		hub.UpdateMaxTransaction()
		for {
			select {
			case client := <-hub.Register:
				hub.UpdateMaxTransaction()
				if len(hub.Transaction) >= hub.MaxTransaction {
					hub.PushWaiting(client)
					continue
				}
				hub.PushTransaction(client)
				client.InformNext(client)
				hub.BroadcastStock()
			case client := <-hub.Unregister:
				hub.WalkRemove(client)
				hub.UpdateMaxTransaction()
				hub.BroadcastNext()
				hub.BroadcastStock()
			case ops := <-hub.Operation:
				switch ops.command {
				case dto.WSOCKET_ENUM_WITH_MERCH_REQUEST:
					ops.client.SetWithMerch(true)
				case dto.WSOCKET_ENUM_NO_MERCH_REQUEST:
					ops.client.SetWithMerch(false)
				}
				hub.BroadcastStock()
			case userID := <-hub.Done:
				hub.ClientDone(userID)
			}
		}
	}()

	return hub
}

func (Hub *queueHub) ClientDone(userID string) {
	Hub.Lock()
	defer Hub.Unlock()

	for _, c := range Hub.Transaction {
		if c.UserID == userID {
			c.Done(nil)
			Hub.Unregister <- c
			return
		}
	}
}

func (Hub *queueHub) UpdateMaxTransaction() {
	Hub.Lock()
	defer Hub.Unlock()

	withMerch, _ := Hub.repository.GetByID(Hub.WithMerchID)
	noMerch, _ := Hub.repository.GetByID(Hub.NoMerchID)
	remainingCapacity := (withMerch.Capacity - withMerch.Registers) + (noMerch.Capacity - noMerch.Registers) + len(Hub.Transaction)

	if remainingCapacity <= Hub.MaxTransaction {
		Hub.MaxTransaction = remainingCapacity
	}
}

func (Hub *queueHub) PushWaiting(client *Client) {
	Hub.Lock()
	defer Hub.Unlock()

	client.SetWaiting()
	Hub.Waiting = append(Hub.Waiting, client)
}

func (Hub *queueHub) PushTransaction(client *Client) {
	Hub.Lock()
	defer Hub.Unlock()

	Hub.Transaction = append(Hub.Transaction, client)
}

func (Hub *queueHub) WalkRemove(client *Client) {
	Hub.Lock()
	defer Hub.Unlock()

	for idx, c := range Hub.Transaction {
		if c == client {
			Hub.Transaction = append(Hub.Transaction[:idx], Hub.Transaction[idx+1:]...)
			return
		}
	}

	for idx, c := range Hub.Waiting {
		if c == client {
			Hub.Waiting = append(Hub.Waiting[:idx], Hub.Waiting[idx+1:]...)
			return
		}
	}
}

func (Hub *queueHub) BroadcastNext() {
	Hub.Lock()
	defer Hub.Unlock()

	if len(Hub.Waiting) == 0 {
		return
	}

	if Hub.MaxTransaction == 0 {
		for _, client := range Hub.Waiting {
			client.Notify(dto.ErrWSMainEventFull.Error())
		}

		return
	}

	// pop next in queue
	next := Hub.Waiting[0]
	if len(Hub.Waiting) == 1 {
		Hub.Waiting = nil
	} else {
		Hub.Waiting = Hub.Waiting[1:]
	}

	// push next transaction
	Hub.Transaction = append(Hub.Transaction, next)

	// broadcast
	next.InformNext(next)
	for _, client := range Hub.Waiting {
		client.InformNext(next)
	}
}

func (Hub *queueHub) BroadcastStock() {
	Hub.Lock()
	defer Hub.Unlock()

	noMerch, err1 := Hub.repository.GetByID(Hub.NoMerchID)
	withMerch, err2 := Hub.repository.GetByID(Hub.WithMerchID)
	if err1 != nil || err2 != nil {
		for _, client := range Hub.Transaction {
			client.Notify(dto.ErrWSCommunicateWithDB.Error())
		}
		return
	}

	for _, client := range Hub.Transaction {
		if client.WithMerch {
			withMerch.Registers++
		} else {
			noMerch.Registers++
		}
	}

	message, _ := json.Marshal(dto.S2CMerchStockInfo{
		WithMerch: withMerch.Capacity - withMerch.Registers,
		NoMerch:   noMerch.Capacity - noMerch.Registers,
	})

	for _, client := range Hub.Transaction {
		client.Notify(string(message))
	}
}

func (Hub *queueHub) IsInQueueByUserID(userID string) bool {
	Hub.Lock()
	defer Hub.Unlock()

	for _, c := range Hub.Transaction {
		if c.UserID == userID {
			return true
		}
	}

	for _, c := range Hub.Waiting {
		if c.UserID == userID {
			return true
		}
	}

	return false
}

func (Hub *queueHub) GetRegisterChannel() chan *Client {
	return Hub.Register
}

func (Hub *queueHub) GetUnregisterChannel() chan *Client {
	return Hub.Unregister
}

func (Hub *queueHub) GetOperationChannel() chan operation {
	return Hub.Operation
}

func (Hub *queueHub) GetWaitingLength() int {
	Hub.Lock()
	defer Hub.Unlock()

	return len(Hub.Waiting)
}
