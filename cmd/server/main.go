// Запускает web-сервер по сбору/обработке метрик.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func ServerFactory(serverType string) (server.IServer, error) {
	if serverType == "http" {
		return server.NewHttpServer(), nil
	}
	if serverType == "grpc" {
		return server.NewGRPCServer(), nil
	}
	return nil, errors.New("не удалось определить тип запускаемого сервера")
}

func main() {

	opts := server.ReadOptions()
	err := logger.InitLogger(opts.LogLevel)
	if err != nil {
		fmt.Printf("Не удалось инициализировать лог, %s \n", err.Error())
	}
	logger.Log.Info(*opts)

	serverType, err := ServerFactory("grpc")
	if err != nil {
		log.Fatal(err.Error())
	}
	srv := server.NewServer(serverType, 40, opts)

	logger.Log.Infof("Start server with options: %v", *opts)
	showBuildInfo("Build version: ", buildVersion)
	showBuildInfo("Build date: ", buildDate)
	showBuildInfo("Build commit: ", buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()

	if err := srv.Run(ctx); err != nil {
		panic(err)
	}
}

func showBuildInfo(caption string, v string) {
	if v == "" {
		logger.Log.Infof("%s: N/A", caption)
	} else {
		logger.Log.Infof("%s: %s", caption, v)

	}
}
