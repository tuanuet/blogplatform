package modules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/internal/interfaces/http/router"
	"github.com/aiagent/boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// HTTPModule provides HTTP server with lifecycle management
var HTTPModule = fx.Module("http",
	fx.Provide(newGinEngine, newHTTPServer),
	fx.Invoke(runServer),
)

// newGinEngine creates the Gin engine with all routes configured
func newGinEngine(p router.Params) *gin.Engine {
	return router.New(p)
}

// newHTTPServer creates the HTTP server
func newHTTPServer(engine *gin.Engine, cfg *config.ServerConfig) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

// runServer starts the HTTP server with lifecycle hooks
func runServer(lc fx.Lifecycle, srv *http.Server, cfg *config.ServerConfig) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", map[string]interface{}{
				"port": cfg.Port,
				"mode": cfg.Mode,
			})
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", err, nil)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server...")
			return srv.Shutdown(ctx)
		},
	})
}
