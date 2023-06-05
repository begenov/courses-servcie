package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/begenov/courses-service/internal/config"
	delivery "github.com/begenov/courses-service/internal/delivery/http"
	"github.com/begenov/courses-service/internal/repository"
	"github.com/begenov/courses-service/internal/server"
	"github.com/begenov/courses-service/internal/service"
	"github.com/begenov/courses-service/pkg/cache"
	"github.com/begenov/courses-service/pkg/database"
)

const (
	path = "./.env"
)

func main() {
	cfg, err := config.Init(path)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.OpenDB(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}

	memCache, err := cache.NewMemoryCache(context.Background(), cfg.Redis)
	if err != nil {
		log.Fatalf("error mem cache init: %v", err)
	}

	repos := repository.NewRepository(db)

	service := service.NewService(repos, memCache, cfg.Redis.Ttl)

	handler := delivery.NewHandler(service)

	srv := server.NewServer(cfg, handler.Init())
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Println("server started ", cfg.Server.Port)

	quit := make(chan os.Signal, 1)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("failed to stop server %v", err)
	}
}
