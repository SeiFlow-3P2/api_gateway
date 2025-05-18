package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SeiFlow-3P2/api_gateway/internal/config"
	"github.com/SeiFlow-3P2/api_gateway/internal/middleware"
	boardProto "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	conf *config.Config
	opts []grpc.DialOption
	mux  *runtime.ServeMux
	srv  *http.Server
}

func CustomHeaderMatcher(key string) (string, bool) {
	lower := strings.ToLower(key)

	switch lower {
	case middleware.UserIDHeader:
		return lower, true
	case "authorization":
		return lower, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func NewApp(config *config.Config) *App {
	return &App{
		conf: config,
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		},
		mux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(CustomHeaderMatcher),
		),
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := boardProto.RegisterBoardServiceHandlerFromEndpoint(
		ctx,
		a.mux,
		a.conf.GetBoardServiceAddr(),
		a.opts,
	); err != nil {
		return err
	}

	authMW := middleware.NewAuthMiddleware( /*client, */ a.conf.GetProtectedRoutes())

	var handler http.Handler = a.mux
	handler = authMW.Handler(handler)

	a.srv = &http.Server{
		Addr:    a.conf.GetServerAddr(),
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", a.conf.GetServerAddr())
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
		a.conf.GetShutdownTimeoutDuration(),
	)
	defer cancel()

	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Could not stop server gracefully: %v", err)
		return a.srv.Close()
	}

	return nil
}
