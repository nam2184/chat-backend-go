package routes

import (
	"log"
	"net/http"

	"github.com/nam2184/mymy/middleware"
	r "github.com/nam2184/mymy/routes/controllers"
	"github.com/nam2184/mymy/routes/controllers/auth/key"
	handler "github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/util"

	"github.com/gorilla/mux"
)

func CreateNormalRouter(prefix string, router *mux.Router, opts *middleware.MiddlewareOptions, logger *util.CustomLogger) {
	db, err := SetupDBX()
	if err != nil {
		log.Fatalf(err.Error())
	}

	cfg := key.Config{
		Address:   "http::examplelocalhost:8000",
		Token:     "token",
		MountPath: "secret",
	}

	vaultStore, err := key.NewVaultStore(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	key.InitializeKeyManager(vaultStore, "secret-key1")

	hopts := handler.NewHandlerConfig(handler.WithLogger(logger),
		handler.WithErrorHandler(middleware.NewErrorHandlerErrorResponder()))

	h := r.NewHandler(db, hopts)
	initRouter := router.PathPrefix(prefix).Subrouter()
	initRouter.Use(mux.MiddlewareFunc(middleware.AttachingHeaders(*opts)),
		mux.MiddlewareFunc(middleware.AttachingCognitoMetadata(*opts)),
	)

	initRouter.Path("/auth").Methods(http.MethodPost).HandlerFunc(h.PostAuthenticated)

	normalRouter := router.PathPrefix(prefix).Subrouter()
	normalRouter.Use(mux.MiddlewareFunc(middleware.AttachingHeaders(*opts)),
		mux.MiddlewareFunc(middleware.ValidatingAuthSchema(*opts)),
		mux.MiddlewareFunc(middleware.AttachingCognitoMetadata(*opts)),
	)

	normalRouter.Path("/user").Methods(http.MethodGet).HandlerFunc(h.GetUser)
	normalRouter.Path("/chats").Methods(http.MethodGet).HandlerFunc(h.GetChats)
	normalRouter.Path("/auth/refresh").Methods(http.MethodGet).HandlerFunc(h.GetRefreshToken)
	normalRouter.Path("/messages/{chatID}").Methods(http.MethodGet).HandlerFunc(h.GetMessages)
	normalRouter.Path("/messages/{chatID}/count").Methods(http.MethodGet).HandlerFunc(h.GetMessagesCount)

	wsRouter := router.PathPrefix(prefix).Subrouter()
	wsRouter.Use(mux.MiddlewareFunc(middleware.AttachingHeaders(*opts)),
		mux.MiddlewareFunc(middleware.ValidatingAuthWSSchema(*opts)),
	)

	wsRouter.HandleFunc("/ws/{chatID}", h.HandleWebsocket)
	wsRouter.HandleFunc("/ws/encrypted/{chatID}", h.HandleEncryptedWebsocket)
}

func CreateTestRouter(prefix string, router *mux.Router, opts *middleware.MiddlewareOptions, logger *util.CustomLogger) {
	db, err := SetupDBX()
	if err != nil {
		log.Fatalf(err.Error())
	}

	hopts := handler.NewHandlerConfig(handler.WithLogger(logger),
		handler.WithErrorHandler(middleware.NewErrorHandlerErrorResponder()))

	h := r.NewHandler(db, hopts)
	initRouter := router.PathPrefix(prefix).Subrouter()
	initRouter.Use(mux.MiddlewareFunc(middleware.AttachingHeaders(*opts)),
		mux.MiddlewareFunc(middleware.AttachingCognitoMetadata(*opts)),
	)

	initRouter.Path("/auth").Methods(http.MethodPost).HandlerFunc(h.PostAuthenticated)

	normalRouter := router.PathPrefix(prefix).Subrouter()
	normalRouter.Use(mux.MiddlewareFunc(middleware.AttachingHeaders(*opts)),
		mux.MiddlewareFunc(middleware.ValidatingAuthSchema(*opts)),
		mux.MiddlewareFunc(middleware.AttachingCognitoMetadata(*opts)),
	)

	normalRouter.Path("/auth/refresh").Methods(http.MethodGet).HandlerFunc(h.GetRefreshToken)
	normalRouter.Path("/messages").Methods(http.MethodGet).HandlerFunc(h.GetMessages)
	normalRouter.Path("/ws").HandlerFunc(h.HandleWebsocket)
}
