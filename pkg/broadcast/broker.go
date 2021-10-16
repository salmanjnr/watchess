package broadcast

import (
	"errors"
	"fmt"
)

type updateBroker struct {
	updates        chan GameUpdate
	clients        map[Client]bool
	registerChan   chan Client
	unregisterChan chan Client
	done           chan struct{}
}

func newBroker() updateBroker {
	return updateBroker{
		updates:        make(chan GameUpdate),
		clients:        make(map[Client]bool),
		registerChan:   make(chan Client),
		unregisterChan: make(chan Client),
		done:           make(chan struct{}),
	}
}

// Wrapping updateBroker for clarity

type gameBroker struct {
	updateBroker
}

type roundBroker struct {
	updateBroker
}

func (b *updateBroker) unregister(client Client) {
	client.Close()
	delete(b.clients, client)
}

func (b *updateBroker) close() {
	b.done <- struct{}{}
}

func (b *updateBroker) sendUpdate(update GameUpdate) {
	for client := range b.clients {
		select {
		case client.Updates() <- update:
		default:
			b.unregister(client)
		}
	}
}

func (b *updateBroker) run(errorChan chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			errorChan <- errors.New(fmt.Sprint(r))
		}
	}()
	for {
		select {
		case client := <-b.registerChan:
			b.clients[client] = true
		case client := <-b.unregisterChan:
			b.unregister(client)
		case update := <-b.updates:
			b.sendUpdate(update)
		case <-b.done:
			// make sure to send all updates to clients before closing
			select {
			case update := <-b.updates:
				b.sendUpdate(update)
			default:
				for client := range b.clients {
					client.Close()
				}
				close(b.done)
				close(b.updates)
				close(b.registerChan)
				close(b.unregisterChan)
				return
			}
		}
	}
}
