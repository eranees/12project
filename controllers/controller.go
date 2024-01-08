package controllers

import (
	"encoding/json"
	"fmt"
	"mysql/db"
	"mysql/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("abcdefghijklmnopqrstuvwxyz")

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	var role_id, userid int
	var username string

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := db.DBcon.Query("SELECT id, username, role_id FROM users WHERE username=$1 AND password=$2", credentials.Username, credentials.Password)

	if err != nil {
		fmt.Println("Wrong")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer result.Close()

	for result.Next() {
		err := result.Scan(&userid, &username, &role_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &models.Claims{
		Userid:   userid,
		Username: username,
		Role:     role_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

// Check marks of particular student
func CheckMarks(w http.ResponseWriter, r *http.Request) {
	var total_marks float64
	claims := extractClaimsFromToken(r)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if claims.Role != 1 {
		w.Write([]byte("You don't have permission to check marks"))
		return
	}
	result, err := db.DBcon.Query("SELECT total_marks FROM marks WHERE student_id=$1", claims.Userid)

	if err != nil {
		fmt.Println("You do not have access of other students")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer result.Close()

	for result.Next() {
		err := result.Scan(&total_marks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s! Your total marks are %f", claims.Username, total_marks)))
}

// Student Profile
func StudentProfile(w http.ResponseWriter, r *http.Request) {
	claims := extractClaimsFromToken(r)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if claims.Role != 1 {
		w.Write([]byte("You don't have permission to check profile details"))
		return
	}
	result, err := db.DBcon.Query("SELECT id, name, role_id, username FROM users WHERE id=$1", claims.Userid)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer result.Close()

	var studentProfile models.StudentProfile

	for result.Next() {
		err := result.Scan(&studentProfile.Userid, &studentProfile.Name, &studentProfile.Role_id, &studentProfile.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(studentProfile)
}

// Add Marks To Student
func AddMarks(w http.ResponseWriter, r *http.Request) {
	claims := extractClaimsFromToken(r)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if claims.Role != 2 {
		w.Write([]byte("You don't have permission"))
		return
	}
	var marks models.Marks
	err := json.NewDecoder(r.Body).Decode(&marks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := db.DBcon.Query("SELECT * FROM users WHERE id=$1", marks.StudentId)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer result.Close()
	insert, err := db.DBcon.Exec("INSERT INTO marks (total_marks, student_id) VALUES ($1, $2)", marks.TotalMarks, marks.StudentId)

	if err != nil {
		json.NewEncoder(w).Encode("Student not available")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := insert.RowsAffected(); err == nil && rowsAffected > 0 {
		json.NewEncoder(w).Encode("Added Successfully")
	} else {
		json.NewEncoder(w).Encode("Not Added Successfully")
	}

}

// Only admin can access
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin Dashboard\n"))
	claims := extractClaimsFromToken(r)
	fmt.Println(claims)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claims.Role != 1 {
		w.Write([]byte("You are not an admin"))
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s! Your role is %d", claims.Username, claims.Role)))
}

// Get claims from token
func extractClaimsFromToken(r *http.Request) *models.Claims {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return nil
	}

	tokenStr := cookie.Value
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil
	}

	return claims
}
