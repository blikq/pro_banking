package pkg

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Admin, Customer, Auditor
type Role struct {
	gorm.Model
	Name  string    `json:"name"`
	Users []User    `gorm:"foreignKey:RoleID" json:"users"`
}

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"-"`
	RoleID   uint   `json:"role_id"`
	Role     Role   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`
}
type Transaction struct {
	gorm.Model
	Amount float64	`json:"amount"`
	From User		`json:"from"`
	To User			`json:"to"`
	Message string	`json:"message"`
	Succeeded bool	`json:"succeeded"`
}

func ConnectDB() (*gorm.DB, error) {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	db.AutoMigrate(&Transaction{})

	return db, nil

}

func CreateTransaction(amount float64, from, to User, message string) (Transaction, error) {
	DB, err := ConnectDB()
	if err != nil {
		return Transaction{}, err
	}

	tx := &Transaction{
		Amount:  amount,
		Message: message,
		From:    from,
		To:      to,
		Succeeded: false,
	}
	var toExistingUser User
	result := DB.Where("email = ?", to.Email).First(&toExistingUser)
	if result.Error != nil {
		return Transaction{}, result.Error
	}
	DB.Create(&tx)
	return *tx, nil
}