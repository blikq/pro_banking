package pkg

import (
	"log"

)

func CreateRole() {
	// Create a new role
	db := ConnectDB();

	if db.Error != nil {
		log.Fatal(db.Error)
	}
}

