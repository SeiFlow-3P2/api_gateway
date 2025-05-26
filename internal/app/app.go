package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/SeiFlow-3P2/api_gateway/internal/config"
	"github.com/SeiFlow-3P2/api_gateway/internal/handler"
	"github.com/SeiFlow-3P2/api_gateway/internal/middleware"
	"github.com/SeiFlow-3P2/api_gateway/internal/util"
	authProto "github.com/SeiFlow-3P2/auth_service/pkg/proto/v1"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type App struct {
	conf     *config.Config
	gwmux    *runtime.ServeMux
	router   *gin.Engine
	dialOpts []grpc.DialOption
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
	router := gin.New()
	// router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(otelgin.Middleware(
		"nota.gateway",
		otelgin.WithSpanNameFormatter(func(c *gin.Context) string {
			return fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		}),
	))

	return &App{
		conf:   config,
		router: router,
		gwmux: runtime.NewServeMux(
			runtime.WithOutgoingHeaderMatcher(util.IsHeaderAllowed),
			runtime.WithMetadata(func(
				ctx context.Context,
				req *http.Request,
			) metadata.MD {
				md := make(metadata.MD)

				if userID, ok := ctx.Value(middleware.UserIDHeader).(string); ok {
					md.Set("x-user-id", userID)
					log.Printf("gRPC Metadata: Forwarding x-user-id: %s", userID)
				}

				return md
			}),

			runtime.WithErrorHandler(func(
				ctx context.Context,
				mux *runtime.ServeMux,
				m runtime.Marshaler,
				w http.ResponseWriter,
				r *http.Request,
				err error,
			) {
				log.Printf("GRPC Gateway Error: %v, Path: %s", err, r.URL.Path)
				runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
			}),
		),
		dialOpts: []grpc.DialOption{
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		},
	}
}

func (a *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// shutdownTracer, err := telemetry.NewTracerProvider(ctx, "nota.gateway", "localhost:4317")
	// if err != nil {
	// 	log.Fatalf("failed to create tracer provider: %v", err)
	// }
	// defer func() {
	// 	if err := shutdownTracer(ctx); err != nil {
	// 		log.Printf("failed to shutdown tracer provider: %v", err)
	// 	}
	// }()

	// shutdownMeter, err := telemetry.NewMeterProvider(ctx, "nota.gateway", "localhost:4317")
	// if err != nil {
	// 	log.Fatalf("failed to create meter provider: %v", err)
	// }
	// defer func() {
	// 	if err := shutdownMeter(ctx); err != nil {
	// 		log.Printf("failed to shutdown meter provider: %v", err)
	// 	}
	// }()

	if err := handler.SetupHandlers(ctx, a.conf, a.gwmux, a.dialOpts); err != nil {
		return fmt.Errorf("failed to setup handlers: %w", err)
	}

	conn, err := grpc.NewClient(a.conf.GetAuthServiceAddr(), a.dialOpts...)
	if err != nil {
		return fmt.Errorf("failed to connect to auth client: %w", err)
	}
	defer conn.Close()
	authClient := authProto.NewAuthServiceClient(conn)

	authMW := middleware.NewAuthMiddleware(
		authClient,
		a.conf.GetProtectedRoutes(),
	)

	a.router.Use(authMW.Handler)

	a.router.Group("/v1/*{grpc_gateway}").Any("", gin.WrapH(a.gwmux))

	srv := &http.Server{
		Addr:    a.conf.GetServerAddr(),
		Handler: a.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	return a.gracefulShutdown(ctx, srv)
}

func (a *App) gracefulShutdown(ctx context.Context, srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Received interrupt signal, shutting down server...")
	case <-ctx.Done():
		log.Println("Parent context cancelled, shutting down server...")
	}

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server gracefully stopped")
	return nil
}
