package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"note_service/app/internal/client/user_client"
	"note_service/app/internal/config"
	"note_service/app/internal/note"
	"note_service/app/internal/note/db"
	"note_service/app/pkg/handlers/metric"
	"note_service/app/pkg/logging"
	"note_service/app/pkg/postgres"

	"note_service/app/pkg/shutdown"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("config initializing")
	cfg := config.GetConfig()

	logger.Println("router initializing")
	router := httprouter.New()

	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	postgresClient, err := postgres.NewClient(context.Background(), cfg.PostgreSQL.Host, cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database, logger)
	if err != nil {
		logger.Fatalf("Error creating PostgreSQL client: %v", err)
	}

	noteStorage := db.NewStorage(postgresClient, logger)
	if err != nil {
		panic(err)
	}

	noteService, err := note.NewService(noteStorage, logger)
	if err != nil {
		panic(err)
	}
	userClient := user_client.NewClient(cfg.UserService.URL, "/me", logger)
	notesHandler := note.Handler{
		Logger:      logger,
		NoteService: noteService,
		UserClient:  userClient,
	}
	notesHandler.Register(router)

	logger.Println("start application")
	start(router, logger, cfg)
}

func start(router http.Handler, logger logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Infof("socket path: %s", socketPath)

		logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

		var err error

		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
