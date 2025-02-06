package pkg

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
	// "io"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

)



func StartService() {
	// ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/api/status", getStatus).Methods("GET")
	router.HandleFunc("/api/create-tx", createTransaction).Methods("POST")

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

	if auth() {
		fmt.Fprint(w, "Authenticated")
	} else {
		fmt.Fprint(w, "Not Authenticated")
	}

	

	// var response map[string]interface{}
	// err = json.NewDecoder(resp.Body).Decode(&response)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(response)

	// var transaction Transaction
	// err := json.NewDecoder(r.Body).Decode(&transaction)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// DB, err := ConnectDB()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// DB.Create(&transaction)

	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(transaction)
}

