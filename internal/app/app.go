package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SeiFlow-3P2/api_gateway/internal/config"
	boardProto "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	addr            string
	opts            []grpc.DialOption
	mux             *runtime.ServeMux
	srv             *http.Server
	shutdownTimeout int
}

func NewApp(config *config.Config) *App {
	return &App{
		addr: config.GetServerAddr(),
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		},
		mux:             runtime.NewServeMux(),
		shutdownTimeout: config.GetShutdownTimeout(),
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := boardProto.RegisterBoardServiceHandlerFromEndpoint(
		ctx,
		a.mux,
		"localhost:9090",
		a.opts,
	); err != nil {
		return err
	}

	a.srv = &http.Server{
		Addr:    a.addr,
		Handler: a.mux,
	}

	serverErrors := make(chan error, 1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", a.addr)
		serverErrors <- a.srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		log.Printf("Got signal: %v", sig)
		return a.gracefulShutdown(ctx)
	}
}

func (a *App) gracefulShutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(
		ctx,
		time.Duration(a.shutdownTimeout)*time.Second,
	)
	defer cancel()

	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Could not stop server gracefully: %v", err)
		return a.srv.Close()
	}

	return nil
}
