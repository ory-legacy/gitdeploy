package sse

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ory-am/gitdeploy/storage"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
)

// Broker is a singleton
var brokerInstance *Broker

type Broker struct {
	// All channels
	channels map[string]*channel
	storage  storage.Storage
}

type channel struct {
	// Create a map of Clients, the keys of the map are the channels
	// over which we can push Messages to attached Clients.  (The values
	// are just booleans and are meaningless.)
	//
	clients map[chan *storage.LogEvent]bool

	// Channel into which new Clients can be pushed
	//
	newClients chan chan *storage.LogEvent

	// Channel into which disconnected Clients should be pushed
	//
	defunctClients chan chan *storage.LogEvent
}

func (b *Broker) OpenChannel(name string) *channel {
	c := &channel{
		make(map[chan *storage.LogEvent]bool),
		make(chan (chan *storage.LogEvent)),
		make(chan (chan *storage.LogEvent)),
	}
	log.Printf("Opening channel %s.", name)
	b.channels[name] = c
	return c
}

func (b *Broker) CloseChannel(channel string) {
	delete(b.channels, channel)
	log.Printf("Closed channel %s.", channel)
}

// Start starts a new go routine. It handles
// the addition & removal of Clients, as well as the broadcasting
// of Messages out to Clients that are currently attached.
//
func (b *Broker) Start(channel string) error {
	c, ok := b.channels[channel]
	if !ok {
		return errors.New("Channel does not exist.")
	}

	go func() {
		for {
			if _, ok := b.channels[channel]; !ok {
				// Channel closed
				return
			}

			// TODO Unhack
			nextMessage, err := b.PullNextMessage(c, channel)
			if err != nil {
				log.Printf("An error occured while pulling a message: %s", err.Error())
				time.Sleep(time.Second)
				continue
			}

			select {
			case s := <-c.newClients:
				// There is a new client attached and we
				// want to start sending them Messages.
				c.clients[s] = true
				log.Println("Added new client.")

			case s := <-c.defunctClients:
				// A client has dettached and we want to
				// stop sending them Messages.
				delete(c.clients, s)
				log.Println("Removed client.")

			case s := <-nextMessage:
				for client := range c.clients {
					client <- s
				}

			case <-timeout(100):
			}
		}
	}()
	return nil
}

func (b *Broker) PullNextMessage(c *channel, name string) (chan *storage.LogEvent, error) {
	leChan := make(chan *storage.LogEvent)
	if len(c.clients) > 0 {
		le, err := b.storage.GetNextUnreadMessage(name)
		b.storage.LogEventIsRead(le)
		if err == mgo.ErrNotFound {
			return nil, nil
		} else if err != nil {
			return nil, err
		} else {
			go func() { leChan <- le }()
		}
	}
	return leChan, nil
}

// This Broker method handles and HTTP request
//
func (b *Broker) EventHandler(w http.ResponseWriter, r *http.Request) {
	// TODO This should be done differently. In some middleware e.g.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	vars := mux.Vars(r)
	app := vars["app"]
	log.Printf("A new listener appeared for channel %s", app)

	c, ok := b.channels[app]
	if !ok {
		log.Printf("Channel %s does not exist.", app)
		http.Error(w, "Channel does not exist!", http.StatusInternalServerError)
		return
	}

	// Make sure that the writer supports flushing.
	//
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel, over which the broker can
	// send this client Messages.
	messageChan := make(chan *storage.LogEvent)

	// Add this client to the map of those that should
	// receive updates
	c.newClients <- messageChan

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	defer b.detachClient(c, messageChan)
	go func() {
		<-notify
		b.detachClient(c, messageChan)
	}()

	// Don't close the connection, instead loop until the
	// client have closed the conncetion, sending Messages
	// and flushing the response each time there is a new
	// message to send along.
	for {
		select {
		case <-w.(http.CloseNotifier).CloseNotify():
			// Done
			log.Println("Finished HTTP request at ", r.URL.Path)
			return
		case e := <-messageChan:
			if len(e.Message) > 0 {
				// Write to the ResponseWriter, `w`.
				fmt.Fprintf(w, "data: %s\n\n", e.Message)
				// log.Printf("Sending data %s", e.Message)

				// Flush the response. This is only possible if
				// the response supports streaming.
				f.Flush()
			}
		}
	}

	// Done.
	log.Printf("Finished HTTP request at %s", r.URL.Path)
}

func (b *Broker) detachClient(c *channel, messageChan chan *storage.LogEvent) {
	c.defunctClients <- messageChan
	log.Println("Detached client.")
}

func New(s storage.Storage) *Broker {
	if brokerInstance != nil {
		return brokerInstance
	}
	brokerInstance = &Broker{make(map[string]*channel), s}
	return brokerInstance
}

func timeout(seconds time.Duration) chan bool {
	t := make(chan bool)
	go func() {
		time.Sleep(seconds * time.Millisecond)
		t <- true
	}()
	return t
}
