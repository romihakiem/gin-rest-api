package main

import (
	"gin-rest-api/config"
	"gin-rest-api/models"
	"log"
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

	err = db.Migrator().DropTable(models.User{}, models.Category{}, models.Post{}, models.Comment{})
	if err != nil {
		log.Fatal("table dropping failed")
	}

	err = db.AutoMigrate(models.User{}, models.Category{}, models.Post{}, models.Comment{})
	if err != nil {
		log.Fatal("migration failed")
	}
}
