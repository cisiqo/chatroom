// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	chatroom *Chatroom

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan map[string]interface{}

	// auth to login
	auth bool

	// client user name
	name string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var m interface{}
		if err := c.conn.ReadJSON(&m); err == nil {
			message := m.(map[string]interface{})
			if message["Type"] == "message" {
				if c.chatroom == nil {
					if !c.auth {
						message["Status"] = "error"
						message["Message"] = "您未授权认证，请登录授权"
					} else {
						message["Status"] = "error"
						message["Message"] = "您未加入直播房间"
					}
					message = addTimeStamp(message)
					c.send <- message
				} else {
					message["Username"] = c.name
					message = addTimeStamp(message)
					c.chatroom.updatedAT = time.Now().Unix()
					c.chatroom.forward <- message
				}
			} else if message["Type"] == "login" {
				if message["Username"] != nil {
					username := message["Username"].(string)
					if username != "" {
						c.name = username
						c.auth = true
						message["Status"] = "ok"
						message["Message"] = "授权认证成功"
					} else {
						message["Status"] = "error"
						message["Message"] = "请添加用户姓名"
					}
					message = addTimeStamp(message)
					c.send <- message
				}
			} else if message["Type"] == "join" {
				var roomID string
				if message["RoomID"] == nil {
					chatRoom := newRoom()
					go chatRoom.run()
					roomID = chatRoom.id
					c.chatroom = chatRoom
					c.hub.rooms[roomID] = chatRoom
					chatRoom.join <- c
					defer func() {
						chatRoom.leave <- c
					}()

					message["Message"] = c.name + "进入直播房间"
					message["Type"] = "message"
					message["RoomID"] = roomID
					message = addTimeStamp(message)
					c.chatroom.updatedAT = time.Now().Unix()
					c.chatroom.forward <- message
				} else {
					roomID = message["RoomID"].(string)
					if roomID != "" {
						if chatRoom, ok := c.hub.rooms[roomID]; ok {
							c.chatroom = chatRoom
							chatRoom.join <- c
							defer func() {
								chatRoom.leave <- c
							}()

							message["Message"] = c.name + "进入直播房间"
							message["Type"] = "message"
							message["RoomID"] = roomID
							message = addTimeStamp(message)
							c.chatroom.updatedAT = time.Now().Unix()
							c.chatroom.forward <- message
						} else {
							message["Status"] = "error"
							message["Message"] = "直播房间ID不存在"
							message = addTimeStamp(message)
							c.send <- message
						}
					} else {
						message["Status"] = "error"
						message["Message"] = "直播房间ID为空"
						message = addTimeStamp(message)
						c.send <- message
					}
				}
			} else if message["Type"] == "leave" {
				var roomID string
				if message["RoomID"] == nil {
					message["Status"] = "error"
					message["Message"] = "直播房间ID为空"
					message = addTimeStamp(message)
					c.send <- message
				} else {
					roomID = message["RoomID"].(string)
					if roomID != "" {
						if chatRoom, ok := c.hub.rooms[roomID]; ok {
							c.chatroom = chatRoom
							chatRoom.leave <- c

							message["Message"] = c.name + "退出直播房间"
							message["Type"] = "message"
							message = addTimeStamp(message)
							c.chatroom.updatedAT = time.Now().Unix()
							c.chatroom.forward <- message
						} else {
							message["Status"] = "error"
							message["Message"] = "直播房间ID不存在"
							message = addTimeStamp(message)
							c.send <- message
						}
					} else {
						message["Status"] = "error"
						message["Message"] = "直播房间ID为空"
						message = addTimeStamp(message)
						c.send <- message
					}
				}
			}
		} else {
			break
		}
	}
}

// add timestamp to message
func addTimeStamp(message map[string]interface{}) map[string]interface{} {
	t := time.Now()
	layout := "2006-01-02 15:04:05"
	message["Time"] = t.Format(layout)
	return message
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan map[string]interface{}, messageBufferSize)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
