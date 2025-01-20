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

func StartService() {
	// Create a new role
	db := ConnectDB()

	if db.Error != nil {
		log.Fatal(db.Error)
	}

	router := mux.NewRouter()

	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))

}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

func login(w http.ResponseWriter, r *http.Request) {
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

	email := r.Header.Get("email")
	if email != user.Email {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	var existingUser User
	db := ConnectDB()
	db.Where("email = ?", user.Email).First(&existingUser)

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

	w.WriteHeader(http.StatusOK)
}

func register(w http.ResponseWriter, r *http.Request) {
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
	db := ConnectDB()
	db.Where("email = ?", user.Email).First(&existingUser)

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

	w.WriteHeader(http.StatusOK)
}







