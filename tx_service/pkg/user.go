package pkg

import (
	"os"
	// "bytes"
	// "encoding/json"
	"log"
	"net/http"
	"io"
	"github.com/joho/godotenv"

)

func auth() bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	// Authenticate 
	url := "http://localhost:8099/api/authenticate-admin"
	

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Set("Authorization", os.Getenv("AUTH_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("client: could not read response body: %s\n", err)
	os.Exit(1)
	}
	log.Printf("client: response body: %s\n", resBody)


	if resp.StatusCode != 200 {
		return false
	}

	return true
}


// req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
// if err != nil {
// 	log.Println(err.Error())
// }
// type User_Temp struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// user := User_Temp{
// 	Email:    os.Getenv("TEST_EMAIL"),
// 	Password:  os.Getenv("TEST_PASSWORD"),
// }


// userJSON, err := json.Marshal(user)

// if err != nil {
// 	log.Println(err.Error())
// }

// req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))

// if err != nil {
// 	log.Println(err.Error())
// }

// client := &http.Client{}
// resp, err := client.Do(req)
// if err != nil {
// 	log.Println(err.Error())
// }

// defer resp.Body.Close()

// resBody, err := io.ReadAll(resp.Body)
// if err != nil {
// 	log.Printf("client: could not read response body: %s\n", err)
// }