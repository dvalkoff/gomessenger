package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dvalkoff/gomessenger/internal/backend/config"
	"github.com/dvalkoff/gomessenger/internal/backend/middleware"
	"github.com/dvalkoff/gomessenger/internal/backend/usecases/chat"
	"github.com/dvalkoff/gomessenger/internal/backend/usecases/messaging"
	"github.com/dvalkoff/gomessenger/internal/backend/usecases/user"
)

const (
	connectionStrEnv = "DB_CONNECTION_STR"
	jwtSecretEnv = "JWT_SECRET"
	frontendUrl = "FRONTEND_URL"
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
	friendsService := user.NewFriendsService(userRepository)
	userController := user.NewUserController(userRegistrationUseCase, findUsersUseCase, friendsService)

	chatRepository := chat.NewChatRepository(db)
	createChatUseCase := chat.NewCreateChatUseCase(chatRepository)
	chatSelectionUseCase := chat.NewChatSelection(chatRepository)
	chatController := chat.NewChatController(createChatUseCase, chatSelectionUseCase)

	messagingRepository := messaging.NewMessagingRepository(db)
	messagingHub := messaging.NewMessagingHub(chatRepository, messagingRepository)
	messagingService := messaging.NewMessagingService(messagingHub, messagingRepository)
	messagingController := messaging.NewMessagingConrtoller(messagingService)

	authProvider := middleware.NewAuthenticationProvider(userRepository, os.Getenv(jwtSecretEnv))
	corsProvider := middleware.NewCorsMiddleware(os.Getenv(frontendUrl))

	go messagingHub.Run()

	httpConfig := config.HttpConfig{
		Port:               8080,
		ReadTimeoutMs:      10000,
		WriteTimeoutMs:     10000,
		ShutdownTimeoutSec: 10,
	}
	server := config.SetUpAndRunServer(
		httpConfig,
		corsProvider,
		authProvider,
		userController.RegisterUser(),
		userController.FindUsers(),
		userController.AddFriend(),
		userController.GetFriends(),
		chatController.CreateChat(),
		chatController.AddUserToChat(),
		chatController.GetChats(),
		messagingController.GetRealtimeUpdates(),
	)

	return gracefulShutdown(
		httpConfig.ShutdownTimeoutSec,
		ctx,
		server,
		db,
		messagingHub,
	)
}

func gracefulShutdown(
	timeout int,
	ctx context.Context,
	httpServer *http.Server,
	db *sql.DB,
	messagingHub messaging.MessagingHub) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		slog.Info("Shutting down server")
		shutdownCtx := context.Background()
		shutdownCtx, cancelShutdown := context.WithTimeout(
			shutdownCtx,
			time.Duration(timeout)*time.Second,
		)
		defer cancelShutdown()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error while shutting down http server", "error", err)
		}
		if err := db.Close(); err != nil {
			slog.Error("Error shutting down database connection pool", "error", err)
		}

		hubShutdownChan := messagingHub.Shutdown()
		<-hubShutdownChan
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
