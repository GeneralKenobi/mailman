package httpgin

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/handler/health"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/handler/mailingentry"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/request"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/internal/email"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer(persistenceCtx persistence.Context, emailer email.Service) *Server {
	server := Server{
		persistenceCtx: persistenceCtx,
		emailer:        emailer,
	}
	server.configure()
	return &server
}

type Server struct {
	persistenceCtx persistence.Context
	emailer        email.Service
	httpServer     *http.Server
}

// Run starts the HTTP server and shuts it down gracefully when ctx is cancelled.
func (server *Server) Run(ctx shutdown.Context) {
	go server.listenAndServe()
	server.shutdownOnContextCancellation(ctx)
}

// configure creates a ready-to-use server and stores it in Server.httpServer.
func (server *Server) configure() {
	httpCfg := config.Get().HttpServer
	ginEngine := server.setupGinEngine()
	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpCfg.Port),
		Handler: ginEngine,
	}
}

// setupGinEngine configures routing, middleware and handlers.
func (server *Server) setupGinEngine() *gin.Engine {
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery(), request.ContextMiddleware, request.LogRequestProcessingMiddleware)

	ginEngine.GET("/health", health.HandlerFunc)

	mailingEntryHandler := mailingentry.NewHandler(server.persistenceCtx, server.emailer)
	ginEngine.POST("/api/messages", mailingEntryHandler.CreateHandlerFunc)
	ginEngine.DELETE("/api/messages/:id", mailingEntryHandler.DeleteHandlerFunc)
	ginEngine.POST("/api/messages/send", mailingEntryHandler.SendMailingIdHandlerFunc)

	return ginEngine
}

func (server *Server) listenAndServe() {
	mdctx.Infof(nil, "Starting HTTP server on address %s", server.httpServer.Addr)
	err := server.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		mdctx.Errorf(nil, "HTTP server exited with error: %v", err)
	} else {
		mdctx.Infof(nil, "HTTP server exited")
	}
}

func (server *Server) shutdownOnContextCancellation(ctx shutdown.Context) {
	defer ctx.Notify()

	<-ctx.Done()
	mdctx.Infof(nil, "Context canceled - shutting down HTTP server")
	serverShutdownCtx, cancel := context.WithTimeout(context.Background(), ctx.Timeout())
	defer cancel()

	err := server.httpServer.Shutdown(serverShutdownCtx)
	if err != nil {
		mdctx.Errorf(nil, "Error shutting down HTTP server: %v", err)
	} else {
		mdctx.Infof(nil, "HTTP server shutdown completed")
	}
}
