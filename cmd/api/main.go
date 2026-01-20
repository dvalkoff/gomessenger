package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dvalkoff/gomessenger/internal/config"
	"github.com/dvalkoff/gomessenger/internal/usecases/chat"
	"github.com/dvalkoff/gomessenger/internal/usecases/user"
)

const (
	connectionStrEnv = "DB_CONNECTION_STR"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	dbConfig := config.DbConfig{
		ConnectionStr: os.Getenv(connectionStrEnv),
	}
	db, err := config.NewDb(dbConfig)
	if err != nil {
		return err
	}
	userRepository := user.NewUserRepository(db)
	userRegistrationUseCase := user.NewUserUserRegistrationUseCase(userRepository)
	findUsersUseCase := user.NewFindUsersUseCase(userRepository)
	userController := user.NewUserController(userRegistrationUseCase, findUsersUseCase)

	chatRepository := chat.NewChatRepository(db)
	hub := chat.NewHub(chatRepository)
	go hub.Run()
	createChatUseCase := chat.NewCreateChatUseCase(chatRepository, hub)
	
	chatController := chat.NewChatController(createChatUseCase)
	messagingController := chat.NewMessagingConrtoller(hub)

	httpConfig := config.HttpConfig{
		Port: 8080,
		ReadTimeoutMs: 10000,
		WriteTimeoutMs: 10000,
		ShutdownTimeoutSec: 10,
	}
    server := config.SetUpAndRunServer(
		httpConfig,
		userController.RegisterUser(),
		userController.FindUsers(),
		chatController.CreateChat(),
		chatController.AddUserToChat(),
		messagingController.GetRealtimeUpdates(),
	)

	return gracefulShutdown(
		httpConfig.ShutdownTimeoutSec,
		ctx,
		server,
		db,
	)
}

func gracefulShutdown(timeout int, ctx context.Context, httpServer *http.Server, db *sql.DB) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
    defer cancel()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		fmt.Fprintf(os.Stdout, "Shutting down server\n")
		shutdownCtx := context.Background()
		shutdownCtx, cancelShutdown := context.WithTimeout(
			shutdownCtx,
			time.Duration(timeout) * time.Second,
		)
		defer cancelShutdown()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down database connection pool: %s\n", err)
		}
	})
	wg.Wait()
	return nil
}

func main() {
	// TODO: add shutdown on db connection loss
    ctx := context.Background()
    if err := run(ctx, os.Stdout, os.Args); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}

