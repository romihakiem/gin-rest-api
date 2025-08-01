package main

import (
	"context"
	"gin-rest-api/config"
	"gin-rest-api/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.DBConnect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	rds, err := config.RedisConnect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	router.GetDocs(r)
	router.GetRoute(r, cfg, db, rds)

	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}

	log.Println("server exiting")
}
