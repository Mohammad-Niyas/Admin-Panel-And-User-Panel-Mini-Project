package auth

import (
	"fmt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var SecretKey = []byte("asdfgbnmkoiuytr")

type User struct {
	Email string `json:"username"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func JwtToken(c *gin.Context, email string, role string) {
	tokenkey, err := CreateToken(email, role)
	if err != nil {
		fmt.Println("Failed to create new token:", err)
		return
	}
	session := sessions.Default(c)
	session.Set(role, tokenkey)
	session.Save()
	check := session.Get(role)
	fmt.Println("Stored Token in Session:", check)
}
func CreateToken(email string, role string) (string, error) {
	claims := User{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenkey, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenkey, nil
}
