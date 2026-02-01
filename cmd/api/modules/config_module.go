package modules

import (
	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/pkg/logger"
	"github.com/aiagent/boilerplate/pkg/validator"
	"go.uber.org/fx"
)

// ConfigModule provides configuration-related dependencies
var ConfigModule = fx.Module("config",
	fx.Provide(
		// Load main config (has unique signature, needs wrapper)
		func() (*config.Config, error) {
			return config.Load("config.yaml")
		},
		// Extract sub-configs from main config
		func(c *config.Config) *config.ServerConfig { return &c.Server },
		func(c *config.Config) *config.DatabaseConfig { return &c.Database },
		func(c *config.Config) *config.RedisConfig { return &c.Redis },
		func(c *config.Config) *config.LoggerConfig { return &c.Logger },
	),
	fx.Invoke(initLogger, initValidator),
)

// initLogger initializes the global logger
func initLogger(cfg *config.LoggerConfig) {
	logger.New(cfg)
	logger.Info("Logger initialized")
}

// initValidator initializes the validator
func initValidator() {
	validator.Init()
}
