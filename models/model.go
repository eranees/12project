package models

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Userid   int    `json:"userid"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

type StudentProfile struct {
	Userid   int    `json:"userid"`
	Name     string `json:"name"`
	Role_id  int    `json:"roleid"`
	Username string `json:"username"`
}

type Marks struct {
	StudentId  int     `json:"studentid"`
	TotalMarks float64 `json:"totalmarks"`
}
