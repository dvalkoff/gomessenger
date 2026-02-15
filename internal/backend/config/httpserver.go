package config

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type MiddlewareFunc func(next http.Handler) http.Handler

func SetUpAndRunServer(
	config HttpConfig,
	corsProvider MiddlewareFunc,
	handleAuth MiddlewareFunc,
	handleWsAuth MiddlewareFunc,
	handleLogIn http.Handler,
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
		handleAuth,
		handleWsAuth,
		handleLogIn,
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
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      handler,
		ReadTimeout:  time.Duration(config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeoutMs) * time.Millisecond,
	}

	go func() {
		slog.Info(fmt.Sprintf("listening on addr %s", httpServer.Addr))
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("error listening and serving", "error", err)
		}
	}()

	return httpServer
}

func NewHandlers(
	corsProvider MiddlewareFunc,
	handleAuth MiddlewareFunc,
	handleWsAuth MiddlewareFunc,
	handleLogIn http.Handler,
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

	mux.Handle("POST /signup", handleRegisterUser)
	mux.Handle("POST /signin", handleLogIn)
	mux.Handle("GET /users/{nickname}", handleAuth(handleFindUsers))
	mux.Handle("POST /users/friends/{friendsNickname}", handleAuth(handleAddFriend))
	mux.Handle("GET /users/friends", handleAuth(handleGetFriends))

	mux.Handle("POST /chats", handleAuth(handleCreateChat))
	mux.Handle("GET /chats", handleAuth(getChats))
	// mux.Handle("PATCH /chats/{chatId}/participants", addUserToChat) // TODO: fix {chatId}
	mux.Handle("GET /messaging", handleWsAuth(getRealtimeUpdates))

	var handler http.Handler = mux
	handler = corsProvider(handler)
	// TODO: panic recoverer handler

	return handler
}
