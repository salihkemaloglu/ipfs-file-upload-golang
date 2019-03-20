package token

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	db "github.com/salihkemaloglu/gignox-main-beta-001/mongodb"
)

//CreateTokenEndpoint user token creation
func CreateTokenEndpoint(user db.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return ""
	}
	return tokenString
}

//ValidateMiddleware token validation
func ValidateMiddleware(authorizationToken string) (string, error) {

	if authorizationToken != "" {
		token, err := jwt.Parse(authorizationToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			return "", err
		}
		if token.Valid {
			return "ok", nil
		}
		return "Invalid authorization token", nil

	}
	return "An authorization header is required", nil

}