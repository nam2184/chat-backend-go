package response

import (
	"time"

	qm "github.com/nam2184/generic-queries"
	"github.com/nam2184/mymy/models/db"
	e "github.com/nam2184/mymy/models/errors"
)

type StandardResponse[T any] struct {
	  Array      []T        `json:"array"`
	  Meta       Meta       `json:"meta"`
}

type QueryResponse[T qm.QueryTypes] struct {
	  Status     string          `json:"status"`
	  Type       string          `json:"type"`
	  ID         string          `json:"id"`
	  Message    string          `json:"message"`
	  Error      e.ErrorType     `json:"error"`
	  Payload    []T             `json:"payload"`
}

func NewQueryRes[T qm.QueryTypes](status string, Type string, id string, message string, payload []T) *QueryResponse[T] {
  return &QueryResponse[T]{ Status: status, Type: Type, ID : id, Message: message, Payload: payload}
}

type AuthResponse struct {
  AccessToken  string  `json:"access_token"`
  RefreshToken string  `json:"refresh_token"`
  Expiry       time.Time  `json:"exp"`
  User         db.User    `json:"user"`
}

func NewAuthResponse(access, refresh string, exp time.Time, user db.User) AuthResponse {
  return AuthResponse{AccessToken: access, RefreshToken: refresh, Expiry: exp, User : user}
}


