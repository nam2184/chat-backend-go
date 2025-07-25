package websockets

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/util"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, opts *options.HandlerOptions, hub *ChatHub) {
	// Get chat ID from the query parameters
	vars := mux.Vars(r)
	chatIDstr := vars["chatID"]

	chatID, err := util.AtoiToUint(chatIDstr)
	if err != nil {
		//opts.Problem.HandleError(middleware.NewError(w, r, err))
		opts.Log.Debug().Msgf("getting id failing for some reaoson, %s", err.Error())
		return
	}

	id, username, err := getIDandUsernameFromToken(r.Context())
	if err != nil {
		opts.Log.Debug().Msgf("getting id from token failing for some reaoson, %s", err.Error())
		//opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	opts.Log.Info().Msgf("Chat ID given in query: %s", chatIDstr)

	// Create a new client and register it in the hub
	client := &Client{conn: conn, chatID: chatID}
	hub.register <- client

	if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		//opts.Problem.HandleError(middleware.NewErrorSocket(client.conn, err))
		hub.unregister <- client
		opts.Log.Debug().Msg("ping failing for some reaoson")
		client.conn.Close() // Close the connection immediately if an error occurs
		return
	}
	go client.readMessages(hub, opts, id, username)
}

func HandleWebSocketEncrypted(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, opts *options.HandlerOptions, hub *ChatHub) {
	// Get chat ID from the query parameters
	vars := mux.Vars(r)
	chatIDstr := vars["chatID"]

	chatID, err := util.AtoiToUint(chatIDstr)
	if err != nil {
		//opts.Problem.HandleError(middleware.NewError(w, r, err))
		opts.Log.Debug().Msgf("getting id failing for some reaoson, %s", err.Error())
		return
	}

	id, username, err := getIDandUsernameFromToken(r.Context())
	if err != nil {
		opts.Log.Debug().Msgf("getting id from token failing for some reaoson, %s", err.Error())
		//opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	opts.Log.Info().Msgf("Chat ID given in query: %s", chatIDstr)

	// Create a new client and register it in the hub
	client := &Client{conn: conn, chatID: chatID}
	hub.register <- client

	if err := client.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		//opts.Problem.HandleError(middleware.NewErrorSocket(client.conn, err))
		hub.unregister <- client
		opts.Log.Debug().Msg("ping failing for some reaoson")
		client.conn.Close() // Close the connection immediately if an error occurs
		return
	}
	go client.readEncryptedMessages(hub, opts, id, username)
}

func getIDandUsernameFromToken(ctx context.Context) (uint, string, error) {
	token, ok := ctx.Value("token").(*key.CustomClaims)
	if !ok {
		return 0, "", fmt.Errorf("No path params found in context")
	}
	id, err := strconv.Atoi(token.Subject)
	if err != nil {
		return 0, "", nil
	}
	return uint(id), token.Username, nil
}
