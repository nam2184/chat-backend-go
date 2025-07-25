package controllers

import (
	"net/http"

	"github.com/nam2184/mymy/routes/controllers/auth/refresh"
	"github.com/nam2184/mymy/routes/controllers/chats"
	"github.com/nam2184/mymy/routes/controllers/messages"
	"github.com/nam2184/mymy/routes/controllers/messages/count"
	"github.com/nam2184/mymy/routes/controllers/user"
)

func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages.GetMessages(w, r, h.db, h.opts)
}

func (h Handler) GetMessagesCount(w http.ResponseWriter, r *http.Request) {
	count.GetMessagesCount(w, r, h.db, h.opts)
}

func (h Handler) GetChats(w http.ResponseWriter, r *http.Request) {
	chats.GetChats(w, r, h.db, h.opts)
}

func (h Handler) GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh.GetRefreshToken(w, r, h.db, h.opts)
}

func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user.GetUser(w, r, h.db, h.opts)
}
