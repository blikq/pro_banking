package pkg

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
)

type Permission struct {
	//gorm.Model
	ID   uint
	Name  string
}

// Admin, Customer, Auditor
type Role struct {
	gorm.Model		
    Name string 	`json:"name"`
    Users []*User 	`gorm:"many2many:user_roles" json:"users"`
}

type User struct {
	gorm.Model
	Email string	`json:"email"`
	Password string `json:"-"` // Don't return password in JSON
	Roles []Role 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:user_roles;" json:"roles"`
}

type Claim struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

func ConnectDB() *gorm.DB {
	godotenv.Load()
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Role{})

	if err != nil {
		log.Fatal(err)
	}
	return db
}