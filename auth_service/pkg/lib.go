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
	Success bool	`json:"success"`
	Token   string	`json:"token"`
}

func StartService() {
	ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/api/status", getStatus).Methods("GET")
	router.HandleFunc("/api/login", login).Methods("POST")
	router.HandleFunc("/api/register", register).Methods("POST")
	router.HandleFunc("/api/authenticate-admin", authenticateAdmin).Methods("POST")

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8099", router))
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

	var existingUser User
	res := DB.Where("email = ?", user.Email).First(&existingUser)

	if res.Error != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	errHash := CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(10000 * time.Minute)

	claims := &Claims{
		Role: existingUser.Role,
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

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "token",
	// 	Value:    tokenString,
	// 	Expires:  expirationTime,
	// 	Path:     "/",
	// 	Domain:   "localhost",
	// 	Secure:   false,
	// 	HttpOnly: true,
	// })
	successResponse := Data{Success: true, Token: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(successResponse)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Write(response)
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

	res := DB.Where("email = ?", user.Email).First(&existingUser)

	if res.Error == nil {
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

	successResponse := Data{Success: true}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}



func authenticateAdmin(w http.ResponseWriter, r *http.Request) {
	ConnectDB()

	tokenString := r.Header.Get("Authorization")
	fmt.Println(tokenString)
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := ParseToken(tokenString)

	if err != nil {
		http.Error(w, "Expired Token", http.StatusUnauthorized)
		return
	}

	var user User
	res := DB.Where("email = ?", claims.StandardClaims.Subject).First(&user)

	if res.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Role.Name != "Admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}