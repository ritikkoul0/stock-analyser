package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stock-analyser/database"
	"stock-analyser/logger"
	"stock-analyser/routers"
	"stock-analyser/utils"

	"syscall"

	"golang.org/x/sync/errgroup"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	log.Println("Starting stock analyser !!")
	utils.UpdateVariables()
	logger.NewLogger("development")
	logger.Infof(
		"Configuration loaded: ServerHost=%s, ServerPort=%s",
		utils.Config.ServerHost, utils.Config.ServerPort,
	)

	// Setup context for graceful shutdown
	appContext, stopSignals := signal.NotifyContext(context.Background(), shutdownSignals...)
	defer stopSignals()

	// Initialize database connection
	logger.Info("Initializing database connection...")
	err := database.InitializeConnection(appContext, utils.Config)
	if err != nil {
		logger.Fatalf("Failed to initialize database connection: %v", err)
	}
	logger.Infof("Database connection initialized successfully")

	// Create an error group for managing goroutines
	errGroup, ctx := errgroup.WithContext(appContext)

	// Start HTTP server
	startHTTPServer(ctx, errGroup, utils.Config)

	// Wait for all goroutines to finish
	if err := errGroup.Wait(); err != nil {
		logger.Fatalf("Application encountered an error: %v", err)
	}

}

func startHTTPServer(ctx context.Context, errGroup *errgroup.Group, config *utils.AppConfig) {
	address := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	logger.Infof("Preparing to start server at %s...", address)

	// Setup router
	router := routers.SetupRouter()

	if router == nil {
		logger.Fatalf("Router setup failed")
	}
	logger.Infof("Router setup completed successfully")

	httpServer := &http.Server{
		Addr:    address,
		Handler: router.Handler(),
	}

	// Start server in a goroutine
	errGroup.Go(func() error {
		logger.Infof("Starting server at %s...", address)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server encountered an error: %v", err)
			return err
		}
		logger.Infof("Server stopped listening for requests")
		return nil
	})

	// Graceful shutdown
	errGroup.Go(func() error {
		<-ctx.Done()
		logger.Infof("Initialising graceful shutdown of server ...")

		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Errorf("Failed to shutdown the server: %s", err)
			return err
		}

		logger.Infof("Server shutdown completed successfully")
		return nil
	})
}
