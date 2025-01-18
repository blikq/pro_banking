package pkg

import(
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
	"log"
	"os"
)

type Permission struct {
	//gorm.Model
	ID   uint
	Name  string
}

// Admin, Customer, Auditor
type Role struct {
	gorm.Model
    Name string `json:"name"`
    Users []*User `gorm:"many2many:user_roles"`
}

type User struct {
	gorm.Model
	Email string
	Password string
	Roles []Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:user_roles;"`
}

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	return db
}