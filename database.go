package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//CreateConnection creates the databas connections using github.com/jinzhu/gorm
func CreateConnection() (*gorm.DB, error) {
	// Get the env variables
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")

	// TODO: handle err?
	return gorm.Open(
		"postgres",
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DBUser, DBPassword, DBHost, DBName),
	)
}
