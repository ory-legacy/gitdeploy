# ory-am/event

[![Build Status](https://travis-ci.org/ory-am/event.svg)](https://travis-ci.org/ory-am/event)

A **very** minimalistic EventManager written in Go. Listeners are being called asynchronous.

## Usage

To install this library (you'll need Golang installed as well), run:
```
$ go get github.com/ory-am/event
```

Simple example:
```go
package main

import "github.com/ory-am/event"
import "fmt"

type eventListener struct {}

func (l *eventListener) Trigger(event string, data interface{}) {
    fmt.Printf("Event %s was triggered with data %s", s)
}

func main() {
    em := event.New()
    em.AttachListener("foobar", l)

    // TriggerAndWait blocks until all listeners are finished
    em.TriggerAndWait("foobar", "data")
}
```

Advanced example using channels:
```go
package main

import "github.com/ory-am/event"
import "fmt"

type eventListener struct {}

const eventName = "foobar"

type eventData struct {
    c chan string
}

func (l *eventListener) Trigger(event string, data interface{}) {
    if d, ok := data.(*eventData); ok {
        d.c <- "it worked!"
        return
    }
    panic("Type assertion failed")
}

func main() {
    data := &eventData{make(chan bool)}
    em := event.New()
    em.AttachListener("foobar", l)

    // Trigger does not wait for listeners to finish.
    // Use channels instead to keep things in sync.
    em.Trigger("foobar", data)

    fmt.Println(<- data.c)
}
```