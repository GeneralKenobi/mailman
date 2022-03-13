package api

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ServeHttp starts an HTTP server based on httpConfig and sets up routing, middleware, handlers. The server shuts down gracefully when ctx
// is cancelled.
func ServeHttp(ctx shutdown.Context, httpConfig config.HttpServer) {
	ginEngine := setupGinEngine()
	// TODO: More configuration
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", httpConfig.Port),
		Handler: ginEngine,
	}

	go runServer(&server)
	shutdownOnContextCancellation(ctx, &server)
}

func setupGinEngine() *gin.Engine {
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery(), requestContextMiddleware, logRequestProcessingMiddleware)

	return ginEngine
}

func runServer(server *http.Server) {
	mdctx.Infof(nil, "Starting HTTP server on address %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		mdctx.Errorf(nil, "HTTP server exited with error: %v", err)
	} else {
		mdctx.Infof(nil, "HTTP server exited")
	}
}

func shutdownOnContextCancellation(ctx shutdown.Context, server *http.Server) {
	<-ctx.Done()
	mdctx.Infof(nil, "Context canceled - shutting down HTTP server")
	serverShutdownCtx, cancel := context.WithTimeout(context.Background(), ctx.Timeout())
	defer cancel()

	err := server.Shutdown(serverShutdownCtx)
	if err != nil {
		mdctx.Errorf(nil, "Error shutting down HTTP server: %v", err)
	} else {
		mdctx.Infof(nil, "HTTP server shutdown completed")
	}
	ctx.Notify()
}
