package token

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Generate(userID int, apiSecret []byte, tokenLifespan time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(tokenLifespan).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(apiSecret)
}

func GetCurrentUserID(c *gin.Context, apiSecret []byte) (int, error) {
	tokenString := ExtractToken(c.Request)
	tkn, err := ParseString(tokenString, apiSecret)
	if err != nil {
		return 0, err
	}
	claims, _ := tkn.Claims.(jwt.MapClaims)
	id, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["user_id"]), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func ParseString(tokenString string, apiSecret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return apiSecret, nil
	})
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
