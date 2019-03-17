package main

import (
	"log"
	"net/http"
	"time"

	"github.com/johannesboyne/sse"
)

func main() {
	sseBroker := sse.NewBroker()

	ticker := time.NewTicker(1000 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				sampleJson := `{"test": "hello world"}`
				// --------------------- USAGE
				// writing to all connections
				sseBroker.Notifier <- []byte(sampleJson)
				// ---------------------/USAGE
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	http.Handle("/events", sseBroker)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
