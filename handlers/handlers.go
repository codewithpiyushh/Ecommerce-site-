package handlers

import (
	"encoding/json"

	"html/template"

	"fmt"
	"net/http"
	"time"

	"fashion/database"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var Con = database.Con // your db connection

var jwtKey = []byte("your_secret_key") 

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func HomelivingHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home_furnishing.html"))
	tmpl.Execute(w, nil)
}

func MenHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/mens.html"))
	tmpl.Execute(w, nil)
}

func WomenHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/women.html"))
	tmpl.Execute(w, nil)
}

func CartHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/cart.html"))
	tmpl.Execute(w, nil)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func Servesignup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/signup.html"))
	tmpl.Execute(w, nil)
}

// ------------------ SIGNUP ------------------

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users(username, email, password,hashedpassword) VALUES(?, ?, ?, ?)"
	result, err := Con.Exec(query, user.Username, user.Email, user.Password, string(hashedPassword))

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	})

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error fetching last insert ID:", err)
	} else {
		fmt.Println("Inserted ID:", id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
	})
}

// ------------------ LOGIN ------------------

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	var creds User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid input",
		})
		return
	}

	var storedHashedPassword string
	err = Con.QueryRow("SELECT hashedpassword FROM users WHERE email = ?", creds.Email).Scan(&storedHashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(creds.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}
	var username string
	err = Con.QueryRow("SELECT username FROM users WHERE email = ?", creds.Email).Scan(&username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Generate JWT
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: creds.Email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Could not generate token",
		})
		return
	}

	// Set token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   tokenStr,
		"message": "Login successful",
	})
}
//------------------Dashboard------------------
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, map[string]string{
		"Username": claims.Username,
		"Email":    claims.Email,
	})
}

// ------------------ JWT MIDDLEWARE ------------------

// func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenStr := r.Header.Get("Authorization")
// 		if tokenStr == "" {
// 			http.Error(w, "Missing token", http.StatusUnauthorized)
// 			return
// 		}

// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Token is valid
// 		next(w, r)
// 	}
// }

func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login.html", http.StatusSeeOther)
				return
			}
			http.Error(w, "Error accessing cookie: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Parse token
		tokenStr := cookie.Value
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Invalid token
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}

		fmt.Println("Authenticated user email:", claims.Email)

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}


func AuthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("token")
	tokenStr := cookie.Value

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not authenticated"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user":          claims.Email,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Clear the session (if you use server-side session storage)
    // Example: sessionStore.Destroy(r.Context(), sessionID) - depends on your implementation

    // 2. Clear the cookie by setting it with a past expiration date
    http.SetCookie(w, &http.Cookie{
        Name:     "token", // or your cookie name
        Value:    "",
        Path:     "/",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
        HttpOnly: true,
    })
	 http.Redirect(w, r, "/index.html", http.StatusSeeOther)
}