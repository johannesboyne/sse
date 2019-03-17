# Simple Server-Side Event Framework

Inspired by [bootic](https://github.com/bootic), thus all credits are directed
there [bootic / broker.go](https://github.com/bootic/bootic_data_collector/blob/master/firehose/broker.go).

## Simple to use

```golang
// setup
sseBroker := sse.NewBroker()
// simply add a handler
http.Handle("/events", sseBroker)
log.Fatal(http.ListenAndServe(":8080", nil))
// alternative
// log.Fatal(http.ListenAndServe(":8080", sseBroker))

// somewhere else
// write to the notification channel
sseBroker.Notifier <- []byte(sampleJson)
// => gets propagated to all connected clients
// clients are identified with ULIDs
```

License: MIT
