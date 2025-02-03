package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

)



func StartService() {
	ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/api/status", getStatus).Methods("GET")
	router.HandleFunc("/api/transaction", createTransaction).Methods("POST")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8092", router))

}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// dsn := os.Getenv("DB_URL")

	// Authenticate 
	url := "http://localhost:8091/api/login"
	type User_Temp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	user := User_Temp{
		Email:    os.Getenv("TEST_EMAIL"),
		Password:  os.Getenv("TEST_PASSWORD"),
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Print(err.Error())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(userJSON))
	if err != nil {
		fmt.Println(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(response)

	var transaction Transaction
	// err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	DB, err := ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	DB.Create(&transaction)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
}