package websockets

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	qm "github.com/nam2184/mymy/models/db"
	"github.com/nam2184/mymy/models/request"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/services/he"
	"github.com/nam2184/mymy/util"
)

type Client struct {
	conn   *websocket.Conn
	chatID uint
}

type ChatHub struct {
	log                *util.CustomLogger
	clients            map[*Client]bool
	register           chan *Client
	unregister         chan *Client
	broadcast          chan qm.Message
	encryptedBroadcast chan qm.EncryptedMessage
	chatClients        map[uint]map[*Client]bool // Maps chat_id to clients in that chat
	userToClient       map[string]*Client        // Maps user_id to their WebSocket connection
	mutex              sync.Mutex
	dbWriter           chan qm.Message
	encryptedDBWriter  chan qm.EncryptedMessage
	heService          *he.HEService
	db                 *sqlx.DB
}

func NewChatHub(db *sqlx.DB, log *util.CustomLogger) *ChatHub {
	return &ChatHub{
		log:         log,
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan qm.Message),
		chatClients: make(map[uint]map[*Client]bool),
		dbWriter:    make(chan qm.Message, 2),
		db:          db,
	}
}

func NewChatHubEncrypted(db *sqlx.DB, log *util.CustomLogger) *ChatHub {
	return &ChatHub{
		log:                log,
		clients:            make(map[*Client]bool),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		encryptedBroadcast: make(chan qm.EncryptedMessage),
		chatClients:        make(map[uint]map[*Client]bool),
		encryptedDBWriter:  make(chan qm.EncryptedMessage, 2),
		db:                 db,
	}
}

func (hub *ChatHub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.addClient(client)

		case client := <-hub.unregister:
			hub.removeClient(client)

		case message := <-hub.broadcast:
			hub.sendMessageToChat(message)

		case message := <-hub.dbWriter:
			hub.handleDBWrites(message)
		}
	}
}

func (hub *ChatHub) RunEncrypted() {
	for {
		select {
		case client := <-hub.register:
			hub.addClient(client)

		case client := <-hub.unregister:
			hub.removeClient(client)

		case message := <-hub.encryptedBroadcast:
			hub.sendEncryptedMessageToChat(message)

		case message := <-hub.encryptedDBWriter:
			hub.handleDBWritesEncrypted(message)
		}
	}
}

// addClient registers a new client in the appropriate chat room
func (hub *ChatHub) addClient(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.log.Info().Msgf("New client added connection in chat ID: %d", client.chatID)
	if hub.chatClients[client.chatID] == nil {
		hub.chatClients[client.chatID] = make(map[*Client]bool)
	}
	hub.chatClients[client.chatID][client] = true
}

// removeClient removes a client from its chat room
func (hub *ChatHub) removeClient(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if _, ok := hub.chatClients[client.chatID][client]; ok {
		delete(hub.chatClients[client.chatID], client)
		hub.closeConn(client.conn)
	}
}

// sendMessageToChat broadcasts a message to all clients in the specified chat room
func (hub *ChatHub) sendMessageToChat(message qm.Message) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.log.Info().Msgf("Sending messages from %d to a connection", message.ChatID)
	for client := range hub.chatClients[message.ChatID] {
		if err := client.conn.WriteJSON(message); err != nil {
			log.Println("Write error:", err)
			hub.unregister <- client // Unregister client on write failure
		}
	}
	log.Println("What da heck its finished:")
}

func (hub *ChatHub) sendEncryptedMessageToChat(message qm.EncryptedMessage) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.log.Info().Msgf("Sending messages from %d to a connection", message.ChatID)
	for client := range hub.chatClients[message.ChatID] {
		if err := client.conn.WriteJSON(message); err != nil {
			log.Println("Write error:", err)
			hub.unregister <- client // Unregister client on write failure
		}
	}
	log.Println("What da heck its finished:")
}

func (client *Client) readMessages(hub *ChatHub, opts *options.HandlerOptions, userID uint, username string) {
	defer func() {
		hub.unregister <- client
	}()

	for {
		var temp request.WSMessage
		var message qm.Message
		// Read a message from the WebSocket
		if err := client.conn.ReadJSON(&temp); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				opts.Log.Debug().Msg("Client Disconnected")
				//opts.Problem.HandleError(middleware.NewErrorSocket(client.conn, err))
				break
			}
			opts.Log.Debug().Msgf("Read error: %s", err)
			return // Exit loop if error occurs (client disconnects or other errors)
		}

		message, err := temp.ConvertToMessageDB()
		if err != nil {
			opts.Log.Info().Msgf("Message Content conversion error: %v", err)
			continue
		}

		// Debug the received message
		// Skip empty messages
		if message.Content == "" && message.Type != "typing" && message.Image == nil {
			opts.Log.Debug().Msgf("No content in message, skipping")
			continue
		}

		if message.Type == "typing" {
			opts.Log.Debug().Msgf("is typing")
			hub.broadcast <- message
			continue
		}

		// Add the chat ID to the message
		message.ChatID = client.chatID
		message.SenderID = userID
		message.SenderName = username
		// Broadcast the message to other clients via the hub
		hub.broadcast <- message
		hub.dbWriter <- message
	}
}

func (client *Client) readEncryptedMessages(hub *ChatHub, opts *options.HandlerOptions, userID uint, username string) {
	defer func() {
		hub.unregister <- client
	}()

	for {
		var temp request.WSEncryptedMessage
		var message qm.EncryptedMessage
		// Read a message from the WebSocket
		if err := client.conn.ReadJSON(&temp); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				opts.Log.Debug().Msg("Client Disconnected")
				//opts.Problem.HandleError(middleware.NewErrorSocket(client.conn, err))
				break
			}
			opts.Log.Debug().Msgf("Read error: %s", err)
			return // Exit loop if error occurs (client disconnects or other errors)
		}

		message, err := temp.ConvertToEncryptedMessageDB()
		if err != nil {
			opts.Log.Info().Msgf("Message Content conversion error: %v", err)
			continue
		}

		// Debug the received message
		// Skip empty messages
		if message.Content == "" && message.Type != "typing" && message.Image == nil {
			opts.Log.Debug().Msgf("No content in message, skipping")
			continue
		}

		// Add the chat ID to the message
		message.ChatID = client.chatID
		message.SenderID = userID
		message.SenderName = username
		// Broadcast the message to other clients via the hub
		hub.encryptedBroadcast <- message
		hub.encryptedDBWriter <- message
	}
}

func (hub *ChatHub) handleDBWrites(message qm.Message) {
	tx, err := hub.db.Beginx()
	if err != nil {
		log.Println("Failed to start transaction", err)
	}

	// Once you have collected messages, write them to the DB
	if err := writeMessageDB(tx, message); err != nil {
		log.Println("Failed to write to DB:", err)
	}
}

func (hub *ChatHub) handleDBWritesEncrypted(message qm.EncryptedMessage) {
	tx, err := hub.db.Beginx()
	if err != nil {
		log.Println("Failed to start transaction", err)
	}

	// Once you have collected messages, write them to the DB
	if err := writeEncryptedMessageDB(tx, message); err != nil {
		log.Println("Failed to write to DB:", err)
	}
}
