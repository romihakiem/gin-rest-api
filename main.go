package main

import (
	"gin-rest-api/config"
	"gin-rest-api/router"
	"log"

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

	log.Printf("Starting HTTP server on port %s", cfg.APPPort)
	if err := r.Run(":" + cfg.APPPort); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
