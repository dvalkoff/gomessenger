package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func SetUpAndRunServer(
	config HttpConfig,
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) *http.Server {
	handler := NewHandlers(
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
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		handleRegisterUser,
		handleFindUsers,
		handleCreateChat,
		addUserToChat,
		getChats,
		getRealtimeUpdates,
	)
	var handler http.Handler = mux
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) {
	mux.Handle("POST /signup", handleRegisterUser)
	// mux.Handle("POST /signin", nil) // TODO
	mux.Handle("GET /users/{nickname}", handleFindUsers)

	mux.Handle("POST /users/{nickname}/chats", handleCreateChat)
	mux.Handle("GET /users/{nickname}/chats", getChats)
	mux.Handle("PATCH /users/{nickname}/chats/{chatId}/participants", addUserToChat)
	
	mux.Handle("GET /users/{nickname}/messaging", getRealtimeUpdates)
}
