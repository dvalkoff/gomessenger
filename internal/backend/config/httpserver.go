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
	corsProvider middleware.CorsMiddleware,
	authProvider middleware.AuthenticationProvider,
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleAddFriend http.Handler,
	handleGetFriends http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) *http.Server {
	handler := NewHandlers(
		corsProvider,
		authProvider,
		handleRegisterUser,
		handleFindUsers,
		handleAddFriend,
		handleGetFriends,
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
	corsProvider middleware.CorsMiddleware,
	authProvider middleware.AuthenticationProvider,
	handleRegisterUser http.Handler,
	handleFindUsers http.Handler,
	handleAddFriend http.Handler,
	handleGetFriends http.Handler,
	handleCreateChat http.Handler,
	addUserToChat http.Handler,
	getChats http.Handler,
	getRealtimeUpdates http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /signup", 					  handleRegisterUser)
	mux.Handle("POST /signin", 					  authProvider.LogIn())
	mux.Handle("GET /users/{nickname}", 		  authProvider.AuthMiddleware(handleFindUsers))
	mux.Handle("POST /users/friends/{friendsNickname}",  authProvider.AuthMiddleware(handleAddFriend))
	mux.Handle("GET /users/friends",  authProvider.AuthMiddleware(handleGetFriends))

	mux.Handle("POST /chats",    authProvider.AuthMiddleware(handleCreateChat))
	mux.Handle("GET /chats", 	  authProvider.AuthMiddleware(getChats))
	// mux.Handle("PATCH /chats/{chatId}/participants", addUserToChat) // TODO: fix {chatId}
	mux.Handle("GET /messaging", authProvider.AuthWsMiddleware(getRealtimeUpdates))

	var handler http.Handler = mux
	handler = corsProvider.HandleCors(handler)
	return handler
}
