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
) *http.Server {
	handler := NewHandlers(handleRegisterUser)
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

func NewHandlers(handleRegisterUser http.Handler) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, handleRegisterUser)
	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux, handleRegisterUser http.Handler) {
	mux.Handle("POST /users", handleRegisterUser)
}
