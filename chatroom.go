package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Chatroom struct {
	id        string
	forward   chan map[string]interface{}
	join      chan *Client
	leave     chan *Client
	clients   map[*Client]bool
	updatedAT int64
}

func newRoom() *Chatroom {
	id := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
	return &Chatroom{
		id:        id,
		forward:   make(chan map[string]interface{}),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		clients:   make(map[*Client]bool),
		updatedAT: time.Now().Unix(),
	}
}

func (c *Chatroom) run() {
	for {
		select {
		case client := <-c.join:
			c.clients[client] = true
		case client := <-c.leave:
			delete(c.clients, client)
		case msg := <-c.forward:
			for target := range c.clients {
				select {
				case target.send <- msg:
					fmt.Println("消息发送成功")
				default:
					fmt.Println("消息发送失败")
					delete(c.clients, target)
				}
			}
		}
	}
}
