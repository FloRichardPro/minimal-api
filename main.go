package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/FloRichardPro/minimal-api/internal/configuration"
	"github.com/FloRichardPro/minimal-api/internal/controllers"
	"github.com/FloRichardPro/minimal-api/internal/routers"
	"github.com/FloRichardPro/minimal-api/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Log *zap.Logger
)

func main() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if err := configuration.Config.LoadConf("/config"); err != nil {
		panic(err)
	}

	if err := services.Init(); err != nil {
		panic(err)
	}

	if err := controllers.Init(); err != nil {
		panic(err)
	}

	if err := routers.Init(); err != nil {
		panic(err)
	}

	routerGin := routers.Router
	addrGin := routers.Config.Addr + ":" + strconv.Itoa(routers.Config.Port)
	srv := &http.Server{
		ReadHeaderTimeout: time.Millisecond,
		Addr:              addrGin,
		Handler:           routerGin,
	}

	go RunGin(addrGin, routerGin)

	WaitSignalShutdown(srv)
}

func WaitSignalShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	Log.Info("Shutdown Server ...")

	// if err := connectors.Close(); err != nil {
	// 	log.Error("error during storage.Close()", zap.Error(err))
	// }

	// Time to wait before force closing
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(routers.Config.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		Log.Error("Server Shutdown: ", zap.Error(err))
	}

	Log.Info("Server exiting")
}

func RunGin(addr string, engine *gin.Engine) {
	Log.Info("REST API listening on : "+addr,
		zap.String("package", "main"))

	Log.Error(engine.Run(addr).Error(),
		zap.String("package", "main"))
}
