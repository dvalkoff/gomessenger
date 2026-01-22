package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dvalkoff/gomessenger/internal/backend/middleware"
)

func SetUpAndRunServer(
	config HttpConfig,
	authProvider middleware.AuthenticationProvider,
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) *http.Server {
	handler := NewHandlers(
		authProvider,
		handleRegisterUser,
		handleFindUsers,
		handleCreateChat,
		addUserToChat,
		getChats,
		getRealtimeUpdates,
	)
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
		Handler: handler,
		ReadTimeout: time.Duration(config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeoutMs) * time.Millisecond,
	}
	
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	return httpServer
}

func NewHandlers(
	authProvider middleware.AuthenticationProvider,
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /signup", 					  handleRegisterUser)
	mux.Handle("POST /signin", 					  authProvider.LogIn()) // TODO
	mux.Handle("GET /users/{nickname}", 		  authProvider.AuthMiddleware(handleFindUsers))

	mux.Handle("POST /chats",    authProvider.AuthMiddleware(handleCreateChat))
	mux.Handle("GET /chats", 	  authProvider.AuthMiddleware(getChats))
	// mux.Handle("PATCH /users/{nickname}/chats/{chatId}/participants", addUserToChat) // TODO: fix {chatId}
	mux.Handle("GET /messaging", authProvider.AuthMiddleware(getRealtimeUpdates))
	return mux
}
