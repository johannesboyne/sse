// Inspired by https://github.com/bootic/bootic_data_collector/blob/master/firehose/broker.go
// Credits to @ https://github.com/bootic
package sse

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

type MessageChan chan []byte

func getULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

type Broker struct {
	Notifier       chan []byte
	activeClients  chan MessageChan
	closingClients chan MessageChan
	clients        map[MessageChan]string
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.activeClients:
			b.clients[s] = getULID()
			log.Printf("Client %s added; #%d clients.", b.clients[s], len(b.clients))
		case s := <-b.closingClients:
			id := b.clients[s]
			delete(b.clients, s)
			log.Printf("Client %s removed; #%d clients.", id, len(b.clients))
		case s := <-b.Notifier:
			// event to broadcast
			for clientMsgChan, _ := range b.clients {
				clientMsgChan <- s
			}
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing.
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(MessageChan)

	// Signal the b that we have a new connection
	b.activeClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		fmt.Println("HERE.")
		b.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		b.closingClients <- messageChan
	}()

	// block waiting or messages broadcast on this connection's messageChan
	for {
		// Write to the ResponseWriter
		fmt.Fprintf(w, "%s\n", <-messageChan)
		// Flush the data inmediatly instead of buffering it for later.
		flusher.Flush()
	}
}

// Broker factory
func NewBroker() (b *Broker) {
	// Instantiate a broker
	b = &Broker{
		Notifier:       make(chan []byte, 1),
		activeClients:  make(chan MessageChan),
		closingClients: make(chan MessageChan),
		clients:        make(map[MessageChan]string),
	}

	// Set it running - listening and broadcasting events
	go b.listen()

	return
}
