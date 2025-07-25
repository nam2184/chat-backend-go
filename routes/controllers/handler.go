package controllers

import (
	"github.com/jmoiron/sqlx"
	opts "github.com/nam2184/mymy/routes/controllers/options"
)

type Handler struct {
    db      *sqlx.DB
    opts    *opts.HandlerOptions
}

func NewHandler(db *sqlx.DB, opts *opts.HandlerOptions) *Handler {
  return &Handler{ db : db, opts : opts}
}

