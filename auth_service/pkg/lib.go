package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Data struct {
	success bool
}

func StartService() {
	ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/api/status", getStatus).Methods("GET")
	router.HandleFunc("/api/login", login).Methods("POST")
	router.HandleFunc("/api/register", register).Methods("POST")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8090", router))

}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

func login(w http.ResponseWriter, r *http.Request) {
	// ConnectDB()

	type User_Temp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user User_Temp
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.Header.Get("email")
	if email != user.Email {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	var existingUser User
	DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	errHash := CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Role: existingUser.Roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	log.Println(tokenString);
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	})

	successResponse := Data{success: true}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successResponse)
}

func register(w http.ResponseWriter, r *http.Request) {
	// ConnectDB()

	type User_ struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user User_
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser User

	DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hash, err := GenerateHashPassword(user.Password)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}
	CreateNormalUser(hash, user.Email)

	successResponse := Data{success: true}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successResponse)
	w.WriteHeader(http.StatusOK)
}






