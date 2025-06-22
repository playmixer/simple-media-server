package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"simple-media-server/internal/adapters/api/rest"
	"simple-media-server/internal/adapters/api/rest/logger"
	"simple-media-server/internal/core/config"
	"simple-media-server/internal/core/converter"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	lgr, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	cnvrt, err := converter.New(ctx, lgr)
	if err != nil {
		log.Fatal(err)
	}

	srv := rest.New(cfg.Rest, lgr, cnvrt)

	log.Println("Server starting...")
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.Info("server stop", zap.Error(err))
		}
	}()
	<-ctx.Done()
	lgr.Info("Stopping...")
	srv.Stop()

	log.Println("Server stopped")
	time.Sleep(time.Second)
}
