package websocket

import (
	"encoding/json"
	"math/rand"
	"sync"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/repository"
)

type (
	QueueHub interface {
		PushWaiting(*Client)
		WalkRemove(*Client)
		UpdateMaxTransaction() (entity.Event, entity.Event)

		BroadcastNext()
		BroadcastStock()

		IsEventHandler(string) bool
		IsInQueueByUserID(userID string) bool
		GetClientInTransactionByUserID(string) *Client

		GetWaitingLength() int
		GetRegisterChannel() chan *Client
		GetUnregisterChannel() chan *Client
		GetOperationChannel() chan operation
		GetBasePrice() (int, int)
	}

	queueHub struct {
		Transaction    []*Client
		Waiting        []*Client
		MaxTransaction int
		MaxWaiting     int

		Register   chan *Client
		Unregister chan *Client
		Operation  chan operation
		Done       chan string

		NoMerchID   string
		WithMerchID string

		NoMerchPrice   int
		WithMerchPrice int

		repository repository.EventRepository

		*sync.Mutex
	}

	operation struct {
		client  *Client
		command string
	}
)

func RunConnHub(repo repository.EventRepository, maxT, maxW int, noMerchID, withMerchID string) QueueHub {
	hub := &queueHub{
		Transaction:    make([]*Client, 0),
		Waiting:        make([]*Client, 0),
		MaxTransaction: maxT,
		MaxWaiting:     maxW,

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
		withMerch, noMerch := hub.UpdateMaxTransaction()
		hub.WithMerchPrice = withMerch.Price
		hub.NoMerchPrice = noMerch.Price

		for {
			select {
			// register takes the client and records it for tracking
			case client := <-hub.Register:
				// when the main event is full
				if hub.MaxTransaction <= 0 {
					client.SetWaiting()
					client.Notify(dto.ErrWSMainEventFull.Error())
					continue
				}

				_, noMerch := hub.UpdateMaxTransaction()
				if (noMerch.Capacity - noMerch.Registers) <= 0 {
					client.SetWithMerch(true)
				}

				// if the main event is full, push the client to waiting list
				if len(hub.Transaction) >= hub.MaxTransaction {
					// if len(hub.Waiting) >= hub.MaxWaiting {
					// 	client.Notify(dto.ErrWSWaitingListFull.Error())
					// 	continue
					// }
					hub.PushWaiting(client)
					continue
				}

				// otherwise, push the client to transaction
				hub.PushTransaction(client)
				client.InformNext(client)
				hub.BroadcastStock()

			// remove the client from the queue
			case client := <-hub.Unregister:
				hub.WalkRemove(client) // remove the client from the queue
				if hub.MaxTransaction > 0 {
					hub.UpdateMaxTransaction() // update the max transaction that could happening at a time
				}
				hub.BroadcastNext()  // broadcast the next client in the waiting list (or notify the main event is full)
				hub.BroadcastStock() // broadcast the remaining stock

			// operation is a command to change the client's transaction type
			// 1. with/without merch selection
			// 2. seat selection
			case ops := <-hub.Operation:
				withMerch, noMerch := hub.UpdateMaxTransaction()
				withMerchRemainder := (withMerch.Capacity - withMerch.Registers)
				noMerchRemainder := (noMerch.Capacity - noMerch.Registers)

				switch ops.command {
				case dto.WSOCKET_ENUM_WITH_MERCH_REQUEST:
					if withMerchRemainder <= 0 {
						ops.client.Notify(dto.ErrWSInvalidCommand.Error())
						continue
					}
					ops.client.SetWithMerch(true)
				case dto.WSOCKET_ENUM_NO_MERCH_REQUEST:
					if noMerchRemainder <= 0 {
						ops.client.Notify(dto.ErrWSInvalidCommand.Error())
						continue
					}
					ops.client.SetWithMerch(false)
				}

				hub.BroadcastStock()

				message, _ := json.Marshal(struct {
					Price int `json:"payment_price"`
				}{
					Price: func() int {
						if ops.client.IsWithMerch() {
							return withMerch.Price + rand.Intn(999)
						}
						return noMerch.Price + rand.Intn(999)
					}(),
				})
				ops.client.Notify(string(message))
			}
		}
	}()

	return hub
}

func (Hub *queueHub) UpdateMaxTransaction() (entity.Event, entity.Event) {
	Hub.Lock()
	defer Hub.Unlock()

	withMerch, _ := Hub.repository.GetByID(Hub.WithMerchID)
	noMerch, _ := Hub.repository.GetByID(Hub.NoMerchID)
	remainingCapacity := (withMerch.Capacity - withMerch.Registers) + (noMerch.Capacity - noMerch.Registers) + len(Hub.Transaction)

	// updates the max transaction that could happening at a time
	// in case of the initial set max transaction is larger
	// than the remaining capacity
	if remainingCapacity <= Hub.MaxTransaction {
		Hub.MaxTransaction = remainingCapacity
	}

	return withMerch, noMerch
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

	// if no one is waiting, nothing to broadcast
	if len(Hub.Waiting) == 0 {
		return
	}

	// if max transaction is 0, i.e. the main event is full
	// broadcast the error message to all waiting clients
	// which in turn also signal the handlers to exit
	// and clear the waiting list
	if Hub.MaxTransaction == 0 {
		for _, client := range Hub.Waiting {
			client.Notify(dto.ErrWSMainEventFull.Error())
		}

		// clear the waiting list,
		// since the client upon exit will also unregister it self from the hub
		// this will prevent this loop block to broadcast to be executed again
		// just to save some resource
		Hub.Waiting = nil
		return
	}

	// pop next in queue to forward to transaction
	next := Hub.Waiting[0]
	if len(Hub.Waiting) == 1 {
		Hub.Waiting = nil
	} else {
		Hub.Waiting = Hub.Waiting[1:]
	}

	// forward the client to transaction
	Hub.Transaction = append(Hub.Transaction, next)

	// broadcast the next client,
	// if the client isn't the next to continue to transaction
	// it'll simply updates and inform the user of current queue number
	// they're at
	next.InformNext(next)
	for _, client := range Hub.Waiting {
		client.InformNext(next)
	}
}

func (Hub *queueHub) BroadcastStock() {
	Hub.Lock()
	defer Hub.Unlock()

	// no more allowed transaction, i.e. main event is full
	// return early
	if Hub.MaxTransaction <= 0 {
		return
	}

	// get the current recorded transaction
	noMerch, err1 := Hub.repository.GetByID(Hub.NoMerchID)
	withMerch, err2 := Hub.repository.GetByID(Hub.WithMerchID)
	if err1 != nil || err2 != nil {
		for _, client := range Hub.Transaction {
			client.Notify(dto.ErrWSCommunicateWithDB.Error())
		}
		return
	}

	for _, client := range Hub.Transaction {
		if client.IsWithMerch() {
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

func (Hub *queueHub) GetClientInTransactionByUserID(userID string) *Client {
	Hub.Lock()
	defer Hub.Unlock()

	for _, c := range Hub.Transaction {
		if c.UserID == userID {
			return c
		}
	}

	return nil
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

func (Hub *queueHub) IsEventHandler(eventID string) bool {
	if Hub.NoMerchID == eventID || Hub.WithMerchID == eventID {
		return true
	}

	return false
}

func (Hub *queueHub) GetBasePrice() (int, int) {
	return Hub.WithMerchPrice, Hub.NoMerchPrice
}
