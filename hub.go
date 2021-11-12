// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	rooms map[string]*Chatroom

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Chatroom),
	}
}

func (h *Hub) run() {
	ticker := time.NewTicker(time.Hour * 2)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case <-ticker.C:
			now := time.Now().Unix()
			if h.rooms != nil {
				for k, v := range h.rooms {
					t := now - v.updatedAT
					num := len(h.rooms[k].clients)
					if t >= 3600 && num == 0 {
						delete(h.rooms, k)
					}
				}
			}
		}
	}
}
