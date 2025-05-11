package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwtSecret = os.Getenv("JWT_SECRET")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type User struct {
	ID       string `gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Validate request
	var user User
	if err := database.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid User or Password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	} else {
		fmt.Printf("User %s logged in successfully\n", req.Username)
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
	})

	fmt.Printf("Segredo usado diretamente: %s\n", jwtSecret)

	// Update the SignedString method to use jwtSecret directly
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		http.Error(w, "Token generate error", http.StatusInternalServerError)
		return
	}

	// Adiciona logs para depuração
	fmt.Printf("Token gerado: %s\n", tokenString)

	resp := LoginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ValidateJWT validates a JWT token and returns whether it is valid.
func ValidateJWT(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Atualiza a validação para usar jwtSecret diretamente
		return []byte(jwtSecret), nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return false
	}

	return token.Valid
}
