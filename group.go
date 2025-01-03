package main

import "log"

// Group manages the clients in a particular group by storing them in a map.
// "register" and "unregister" are both channels, that adds and removes a client from the
// the group. It has a broadcast channel that sends the message to every client in the group, using
// the client's send channel.
type Group struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// NewGroup creates a new group using id, registers the group in gh using gh's register channel
// and returns the newly created group.
func NewGroup(id string, gh *GroupsManager) *Group {
	group := &Group{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
	gh.register <- group

	return group
}

// Every group calls it's Run method in one and only one goroutine. Run manages g's clients i.e
/*
 - It listens for messages in g.register channel to add a client to g.clients
 - It listens for messages in g.unregister channel to remove a client from g.clients
 - It listens for messages in g.broadcast channel to send messages to the client's send channel
*/
func (g *Group) Run() {
	defer func() {
		log.Printf("cleared the group's \"Manage\" goroutine with id: %s", g.id)
	}()

	for {
		select {
		case message, ok := <-g.broadcast:
			if !ok {
				for client := range g.clients {
					delete(g.clients, client)
					// when g's broadcast channel is closed, the line below closes every clients' send channel
					close(client.send)
				}
				// After closing all clients' send channel, the line below exits the Run method. This frees up resources and avoids goroutine leak.
				return
			}

			for client := range g.clients {
				// broadcasts the message to every client in g.clients.
				client.send <- message
			}
		case client := <-g.register:
			g.clients[client] = true
		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				// unregistering a client closes the client's send channel.
				close(client.send)
			}
			if len(g.clients) == 0 {
				// when the last client in g unregisters, the line below exits the Run method. This frees up resources and avoids goroutine leak.
				return
			}
		}

	}
}
