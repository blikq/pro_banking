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
	ID   uint
	Name string
}

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

type Claims struct {
	Role Role `json:"role"`
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
	var role Role
	DB.Where("name = ?", "customer").First(&role)
	if role.ID == 0 {
		role = Role{Name: "customer"}
		DB.Create(&role)
	}
	user := User{Email: email, Password: string(hashedPassword), Role: role}
	DB.Create(&user)
	return user
}

func CreateAdminUser(hash, email string) User {
	hashedPassword := hash
	var role Role
	DB.Where("name = ?", "admin").First(&role)
	if role.ID == 0 {
		role = Role{Name: "admin"}
		DB.Create(&role)
	}
	user := User{Email: email, Password: string(hashedPassword), Role: role}
	DB.Create(&user)
	return user
}

func CreateAuditorUser(hash, email string) User {
	hashedPassword := hash
	var role Role
	DB.Where("name = ?", "auditor").First(&role)
	if role.ID == 0 {
		role = Role{Name: "auditor"}
		DB.Create(&role)
	}
	user := User{Email: email, Password: string(hashedPassword), Role: role}
	DB.Create(&user)
	return user
}

