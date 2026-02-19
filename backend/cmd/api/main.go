package main

import (
	"context"
	"database/sql"
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
	"github.com/dvalkoff/gomessenger/internal/backend/usecases/events"
	"github.com/dvalkoff/gomessenger/internal/backend/usecases/user"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	appConfig, err := config.MakeAppConfig()
	if err != nil {
		return err
	}
	err = middleware.InitLogger(w, appConfig.LoggingConfig)
	if err != nil {
		return err
	}
	db, err := config.NewDB(appConfig.DbConfig)
	if err != nil {
		return err
	}
	userRepository := user.NewUserRepository(db)
	pubKeyRepository := user.NewPubKeyRepository(db)
	userService := user.NewUserService(userRepository, pubKeyRepository)
	contactsService := user.NewContactsService(userRepository)
	userController := user.NewUserController(userService, contactsService)

	chatRepository := chat.NewChatRepository(db)
	chatService := chat.NewChatService(chatRepository)
	chatController := chat.NewChatController(chatService)

	eventsRepository := events.NewEventsRepository(db)
	eventsHub := events.NewEventsHub()
	eventsService := events.NewEventService(chatRepository, eventsRepository, eventsHub)
	eventsController := events.NewEventsController(eventsService)

	authProvider := middleware.NewAuthenticationProvider(userRepository, appConfig.HttpConfig.JwtSecret)
	corsProvider := middleware.NewCorsMiddleware(appConfig.HttpConfig.CorsAllowedURL)

	go eventsHub.Run()

	server := config.SetUpAndRunServer(
		appConfig.HttpConfig,
		config.MiddlewareFunc(corsProvider.HandleCors),
		config.MiddlewareFunc(authProvider.AuthMiddleware),
		config.MiddlewareFunc(authProvider.AuthWsMiddleware),
		authProvider.LogIn(),

		userController.RegisterUser(),
		userController.FindUsers(),
		userController.AddContact(),
		userController.GetContacts(),
		
		chatController.CreateChat(),
		chatController.AddUserToChat(),
		chatController.GetChats(),


		eventsController.HandleEventsWS(),
		eventsController.CreateEvent(),
	)

	return gracefulShutdown(
		appConfig.GracefulShutdownConfig,
		ctx,
		server,
		db,
	)
}

func gracefulShutdown(
	config config.GracefulShutdownConfig,
	ctx context.Context,
	httpServer *http.Server,
	db *sql.DB,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		slog.Info("Shutting down server")
		shutdownCtx := context.Background()
		shutdownCtx, cancelShutdown := context.WithTimeout(
			shutdownCtx,
			time.Duration(config.ShutdownTimeoutSec)*time.Second,
		)
		defer cancelShutdown()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error while shutting down http server", "error", err)
		}

		if err := db.Close(); err != nil {
			slog.Error("Error shutting down database connection pool", "error", err)
		}
	})
	wg.Wait()
	return nil
}

func main() {
	// TODO: add shutdown on db connection loss
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		slog.Error("Application was shut down due to an error", "error", err)
		os.Exit(1)
	}
}
