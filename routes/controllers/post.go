package controllers

import (
	"net/http"

	"github.com/nam2184/mymy/routes/controllers/auth"
	"github.com/nam2184/mymy/routes/controllers/auth/signup"
)

func (h Handler) PostAuthenticated(w http.ResponseWriter, r *http.Request) {
	auth.PostAuthenticated(w, r, h.db, h.opts)
}

func (h Handler) PostAuthSignUp(w http.ResponseWriter, r *http.Request) {
	signup.PostAuthSignUp(w, r, h.db, h.opts)
}
