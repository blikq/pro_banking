package pkg

import (
	"log"
	"os"
	"sync"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"

)

var (
    DB   *gorm.DB
    once sync.Once
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



type Claims struct {
	Role []Role `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")



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

func CreateNormalUser(hash, email string) User {
	hashedPassword := hash
	user := User{Email: email, Password: string(hashedPassword), Roles: []Role{{Name: "Customer"}}}
	DB.Create(&user)
	return user
}

func CreateAdminUser(hash, email string) User {
	hashedPassword := hash
	user := User{Email: email, Password: string(hashedPassword), Roles: []Role{{Name: "Admin"}}}
	DB.Create(&user)
	return user
}

func CreateAuditorUser(hash, email string) User {
	hashedPassword := hash
	user := User{Email: email, Password: string(hashedPassword), Roles: []Role{{Name: "Auditor"}}}
	DB.Create(&user)
	return user
}

func CreateTransaction(amount float64, from, to User, message string) (Transaction, error) {
	tx := &Transaction{
		Amount:  amount,
		Message: message,
		From:    from,
		To:      to,
	}
	var toExistingUser User
	result := DB.Where("id = ?", from.ID).First(&toExistingUser)
	if result.Error != nil {
		return Transaction{}, result.Error
	}
	DB.Create(&tx)
	return *tx, nil
}