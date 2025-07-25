package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nam2184/mymy/routes/controllers/websockets"
)

var oncePlain sync.Once
var onceEncrypted sync.Once
var plainHub *websockets.ChatHub
var encryptedHub *websockets.ChatHub

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h Handler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	oncePlain.Do(func() {
		plainHub = websockets.NewChatHub(h.db, h.opts.Log)
		go plainHub.Run()
	})
	websockets.HandleWebSocket(w, r, conn, h.opts, plainHub)
}

func (h Handler) HandleEncryptedWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	onceEncrypted.Do(func() {
		encryptedHub = websockets.NewChatHubEncrypted(h.db, h.opts.Log)
		go encryptedHub.RunEncrypted()
	})
	websockets.HandleWebSocketEncrypted(w, r, conn, h.opts, encryptedHub)
}
