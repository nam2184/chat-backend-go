package main

import (
	"github.com/nam2184/mymy/middleware"
	"github.com/nam2184/mymy/routes"
	"github.com/nam2184/mymy/util"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)


func ServeHTTP() {  
    r := mux.NewRouter()
    opts := middleware.CreateMiddlewareOptions()
    opts.WithErrorHandlerFunc(middleware.NewErrorHandlerErrorResponder())
    
    logger, err := util.NewCustomLogger(true); if err != nil {
      log.Fatal()
    }
    opts.WithLogger(logger)
    routes.CreateNormalRouter("/api/v1", r, opts, logger) 
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Allow requests from this origin
        AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Allow custom headers
        AllowCredentials: true,
        Debug: true, // Optional: For debugging purposes
    })

    handler := c.Handler(r)
    fmt.Printf("Listening \n")
    log.Fatal(http.ListenAndServe(":8000", handler)) 
    defer logger.Close()
}

func main() {
  ServeHTTP()
}
