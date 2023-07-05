package routers

import (
	"fmt"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	Router           *gin.Engine
	ValidateInstance *validator.Validate

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

	ValidateInstance = validator.New()

	Router = gin.New()

	Router.Use(ginzap.RecoveryWithZap(Log, true))
	Router.Use(ginzap.Ginzap(Log, time.RFC3339, true))

	Router.GET("/foos", GetInstanceFooRouter().GetAll)
	Router.Group("/foo").
		POST("", GetInstanceFooRouter().Post).
		GET("/:foo_uuid", GetInstanceFooRouter().Get).
		PUT("/:foo_uuid", GetInstanceFooRouter().Put).
		PATCH("/:foo_uuid", GetInstanceFooRouter().Patch)

	return nil
}
