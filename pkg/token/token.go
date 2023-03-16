package token

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Generate(userId int) (string, error) {

	//todo the token lifespan should get there via paramter. The env variable should be listed somewhere with other env variables
	tokenLifespan, err := strconv.Atoi("1")

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//todo API_SECRET should get there via paramter. The env variable should be listed somewhere with other env variables
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func Valid(c *gin.Context) error {
	tokenString := ExtractToken(c.Request)
	_, err := ParseString(tokenString)

	if err != nil {
		return err
	}
	return nil
}

func ParseString(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//todo API_SECRET should get there via paramter. The env variable should be listed somewhere with other env variables
		return []byte(os.Getenv("API_SECRET")), nil
	})
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
