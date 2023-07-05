package routers

import (
	"fmt"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Router *gin.Engine

	Config Conf
	Log    *zap.Logger
)

// Conf for the routers package.
type Conf struct {
	GinMode         string `mapstructure:"gin_mode"`
	Addr            string `mapstructure:"addr"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

func Init() (err error) {
	Log, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("can't build new logger : %w", err)
	}

	Router = gin.New()

	Router.Use(ginzap.RecoveryWithZap(Log, true))
	Router.Use(ginzap.Ginzap(Log, time.RFC3339, true))

	Router.Group("/foo").
		GET("/:foo_uuid", GetInstanceFooRouter().Get).
		POST("/:foo_uuid", GetInstanceFooRouter().Post).
		PUT("", GetInstanceFooRouter().Put).
		PATCH("", GetInstanceFooRouter().Patch)

	return nil
}
