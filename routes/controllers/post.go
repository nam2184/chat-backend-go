package controllers

import (
	"net/http"

	"github.com/nam2184/mymy/routes/controllers/auth"
)


func (h Handler) PostAuthenticated(w http.ResponseWriter, r *http.Request) {
    auth.PostAuthenticated(w, r, h.db, h.opts)
}
