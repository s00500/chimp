package main

import "encoding/json"

// Broker handles broadcasting messages to clients
type Broker struct {
	clients     map[chan []byte]struct{}
	subscribe   chan chan []byte
	unsubscribe chan chan []byte
	broadcast   chan []byte
}

// NewBroker creates a new Broker
func NewRunningBroker() *Broker {
	b := &Broker{
		clients:     make(map[chan []byte]struct{}),
		subscribe:   make(chan chan []byte),
		unsubscribe: make(chan chan []byte),
		broadcast:   make(chan []byte),
	}
	go b.run()
	return b
}

// Run starts the broker to handle subscriptions, unsubscriptions, and broadcasting
func (b *Broker) run() {
	for {
		select {
		case client := <-b.subscribe:
			b.clients[client] = struct{}{}
		case client := <-b.unsubscribe:
			delete(b.clients, client)
			close(client)
		case message := <-b.broadcast:
			for client := range b.clients {
				// send message to client in a non-blocking manner
				select {
				case client <- message:
				default:
					// if client channel is blocked, skip sending
				}
			}
		}
	}
}

// Subscribe adds a new client to the broker
func (b *Broker) Subscribe() chan []byte {
	client := make(chan []byte, 100) // Buffered channel to avoid blocking
	b.subscribe <- client
	return client
}

// Unsubscribe removes a client from the broker
func (b *Broker) Unsubscribe(client chan []byte) {
	b.unsubscribe <- client
}

type WSMessage struct {
	DataType string
	Data     any
}

// Broadcast sends a message to all subscribed clients
func (b *Broker) Broadcast(mtype string, data interface{}) error {
	jsonData, err := json.Marshal(WSMessage{DataType: mtype, Data: data})
	if err != nil {
		return err
	}

	b.broadcast <- jsonData
	return nil
}

// Broadcast sends a message to all subscribed clients
func (b *Broker) BroadcastString(data string) error {
	//jsonData, err := json.Marshal(WSMessage{DataType: mtype, Data: data})
	//if err != nil {
	//	return err
	//}

	b.broadcast <- []byte(data)
	return nil
}

// Broadcast sends a message to all subscribed clients
func (b *Broker) BroadcastRaw(message []byte) {
	b.broadcast <- message
}
