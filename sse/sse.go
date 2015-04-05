// Golang HTML5 Server Side Events Server based on
// https://github.com/kljensen/golang-html5-sse-example
package sse

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
	"time"
)

// A single Broker will be created in this program. It is responsible
// for keeping a list of which Clients (browsers) are currently attached
// and broadcasting events (Messages) to those Clients.
//
type Broker struct {
	// All channels
	Channels map[string]*Channel
}

type Channel struct {
	// Create a map of Clients, the keys of the map are the channels
	// over which we can push Messages to attached Clients.  (The values
	// are just booleans and are meaningless.)
	//
	Clients map[chan string]bool

	// Channel into which new Clients can be pushed
	//
	NewClients chan chan string

	// Channel into which disconnected Clients should be pushed
	//
	DefunctClients chan chan string

	// Channel into which Messages are pushed to be broadcast out
	// to attahed Clients.
	//
	Messages chan string
}

func (b *Broker) OpenChannel(channel string) *Channel {
	b.Channels[channel] = &Channel{
        make(map[chan string]bool),
        make(chan (chan string)),
        make(chan (chan string)),
        make(chan string),
    }
	log.Printf("Opened channel %s.", channel)
	return b.Channels[channel]
}

func (b *Broker) CloseChannel(channel string) {
	delete(b.Channels, channel)
	log.Printf("Closed channel %s.", channel)
}

func (b *Broker) AddMessage(channel string, message string) {
	c, ok := b.Channels[channel]
	if !ok {
		log.Printf("Channel %s not found, creating...", channel)
		c = b.OpenChannel(channel)
	}
	c.Messages <- message
	log.Printf("Added message %s to channel %s", message, channel)
}

// This Broker method starts a new goroutine.  It handles
// the addition & removal of Clients, as well as the broadcasting
// of Messages out to Clients that are currently attached.
//
func (b *Broker) Start() {
	// Start a goroutine
	//
	go func() {

		// Loop endlessly
		//
		for {
			var wg sync.WaitGroup

			// Loop through all channels
			for channelName, c := range b.Channels {
				log.Printf("Cycling channel %s", channelName)
				wg.Add(1)
				timeout := make(chan bool, 1)
				go func() {
					time.Sleep(5 * time.Second)
					timeout <- true
				}()

				// Handle each channel in a separate go routine
				go func() {
					defer wg.Done()

					// Block until we receive from one of the
					// three following channels.
					select {

					case s := <-c.NewClients:

						// There is a new client attached and we
						// want to start sending them Messages.
						c.Clients[s] = true
						log.Println("Added new client.")

					case s := <-c.DefunctClients:

						// A client has dettached and we want to
						// stop sending them Messages.
						delete(c.Clients, s)
						log.Println("Removed client.")

					case msg := <-c.Messages:

						// TODO unhack...
						log.Println("Waiting for clients.")
						if len(c.Clients) < 1 {
							select {
							case s := <-c.NewClients:
								c.Clients[s] = true
								log.Println("Added new client.")
							case <- timeout:
								log.Println("No client available, moving on.")
								return
							}
						}

						// There is a new message to send.  For each
						// attached client, push the new message
						// into the client's message channel.
						for s := range c.Clients {
							s <- msg
						}
						log.Printf("Broadcast message %s on channel %s to %d clients", msg, channelName, len(c.Clients))
					case <- timeout:
						log.Println("Warning: Nothing is happening, moving on.")
						return
					}
					log.Printf("Ending channel %s's cycle", channelName)
				}()
			}
			wg.Wait()
		}
	}()
}

// This Broker method handles and HTTP request
//
func (b *Broker) EventHandler(w http.ResponseWriter, r *http.Request) {
	// TODO This should be done differently. In some middleware e.g.
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	vars := mux.Vars(r)
	app := vars["app"]
	log.Printf("A new listener appeared for channel %s", app)
	c, ok := b.Channels[app]
	if !ok {
		log.Printf("Channel %s not found, creating...", app)
		c = b.OpenChannel(app)
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
	messageChan := make(chan string)

	// Add this client to the map of those that should
	// receive updates
	c.NewClients <- messageChan

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached Clients
		// when `EventHandler` exits.
		c.DefunctClients <- messageChan
		log.Println("HTTP connection just closed.")
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
		case msg := <-messageChan:
			// Write to the ResponseWriter, `w`.
			fmt.Fprintf(w, "data: %s\n\n", msg)

			// Flush the response. This is only possible if
			// the response supports streaming.
			f.Flush()
		}
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}

func New() *Broker {
	return &Broker{
		make(map[string]*Channel),
	}
}
