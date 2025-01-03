package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Handler provides manager to the http handlers
type Handler struct {
	manager *GroupsManager
}

func NewHandler() *Handler {
	return &Handler{
		manager: NewGroupsManager(),
	}
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if _, ok := h.manager.Contains(id); ok {
		http.Error(w, "id already exists", http.StatusBadRequest)
		return
	}

	h.manager.groups[id] = NewGroup(id, h.manager)

	w.WriteHeader(http.StatusCreated)
	message := map[string]string{
		"message": fmt.Sprintf("created group with id: %s", id),
	}
	res, _ := json.Marshal(&message)
	w.Write(res)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	group, ok := h.manager.Contains(id)
	if !ok {
		http.Error(w, fmt.Sprintf("found no group with the id: %s", id), http.StatusBadRequest)
		return
	}
	h.manager.unregister <- group
	message := map[string]string{
		"message": fmt.Sprintf("deleted group with id: %s", id),
	}
	res, _ := json.Marshal(&message)
	w.Write(res)
}

// ServeWS handles websocket connections. It receives the group id as a path parameter and uses the
// provided id to find the group. It creates a new Client, registers it in the found group and calls
// the clients Read and Write method in goroutines to handle reading and writing to the connection.
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	group, ok := h.manager.Contains(id)
	if !ok {
		http.Error(w, fmt.Sprintf("found no group with the id: %s", id), http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by returning true
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// The group's Run method is called if this new connection is the first connection to the group.
	// This is because there should only be one Run method for every connection "Client" to manage
	// all the Clients of that particular group. Check out the Run method in group.dart for more
	// details about this.
	if len(group.clients) == 0 {
		go group.Run()
	}
	client := NewClient(conn, group)
	group.register <- client

	// Reading from a connection is done here and here alone because there can only be one
	// goroutine per-connection to read from a connection.
	go client.Read()
	// Writing to a connection is done here and here alone because there can only be one
	// goroutine per-connection to write write to a connection.
	go client.Write()
}
