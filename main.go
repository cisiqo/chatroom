// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func creatRoom(hub *Hub, w http.ResponseWriter, r *http.Request) {
	result := make(map[string]interface{})
	chatRoom := newRoom()
	go chatRoom.run()
	roomID := chatRoom.id
	hub.rooms[roomID] = chatRoom
	result["RoomID"] = roomID
	mjson, _ := json.Marshal(result)
	mString := string(mjson)
	fmt.Fprint(w, mString)
}

func getRoomsNumber(hub *Hub, w http.ResponseWriter, r *http.Request) {
	result := make(map[string]interface{})
	result["Rooms_number"] = len(hub.rooms)
	mjson, _ := json.Marshal(result)
	mString := string(mjson)
	fmt.Fprint(w, mString)
}

func getRooms(hub *Hub, w http.ResponseWriter, r *http.Request) {
	result := make(map[string]interface{})
	var rooms []string
	for k := range hub.rooms {
		rooms = append(rooms, k)
	}
	result["Rooms"] = rooms
	mjson, _ := json.Marshal(result)
	mString := string(mjson)
	fmt.Fprint(w, mString)
}

func getRoomMembersNumber(hub *Hub, w http.ResponseWriter, r *http.Request) {
	result := make(map[string]interface{})
	for k := range hub.rooms {
		result[k] = len(hub.rooms[k].clients)
	}
	mjson, _ := json.Marshal(result)
	mString := string(mjson)
	fmt.Fprint(w, mString)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/api/create_room", func(w http.ResponseWriter, r *http.Request) {
		creatRoom(hub, w, r)
	})
	http.HandleFunc("/api/get_rooms_number", func(w http.ResponseWriter, r *http.Request) {
		getRoomsNumber(hub, w, r)
	})
	http.HandleFunc("/api/get_rooms", func(w http.ResponseWriter, r *http.Request) {
		getRooms(hub, w, r)
	})
	http.HandleFunc("/api/get_room_members_number", func(w http.ResponseWriter, r *http.Request) {
		getRoomMembersNumber(hub, w, r)
	})

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
