package pkg

import (
	"log"
)

func CreateRole() {
	// Create a new role
	db := ConnectDB();
	role := Role{Name: "Admin"}
	db.Create(&role)
	if db.Error != nil {
		log.Fatal(db.Error)
	}
}