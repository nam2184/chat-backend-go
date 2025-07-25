package middleware

import (
	"net/http"

	"github.com/nam2184/mymy/util"
)

type AttachHeaders struct {
  next http.Handler
  problem ErrorHandler
  log     *util.CustomLogger
}

func (a *AttachHeaders) ServeHTTP (w http.ResponseWriter, r *http.Request)  {
  a.log.Info().Msg("Attaching CORS headers")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
  
  a.log.Info().Msg("Attaching control allow credentials")
  w.Header().Set("Access-Control-Allow-Credentials", "true")
 
  // Call the next handler
  a.next.ServeHTTP(w, r)
}

