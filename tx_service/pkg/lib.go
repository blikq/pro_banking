package pkg

import (
	"log"
	"os"
	"sync"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount float64	`json:"amount"`
	From User		`json:"from"`
	To User			`json:"to"`
	Message string	`json:"message"`
}

func ConnectDB() {
	once.Do(func() {
        err := godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file")
        }

        dsn := os.Getenv("DB_URL")
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }

        db.AutoMigrate(&User{}, &Role{})

        DB = db
    })

}