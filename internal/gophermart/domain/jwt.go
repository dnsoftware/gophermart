package domain

import (
	"fmt"
	"github.com/dnsoftware/gophermart/internal/constants"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
// передаем ID пользователя
func BuildJWTString(userId int64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.TOKEN_EXP)),
		},

		// собственное утверждение
		UserID: userId,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(constants.SECRET_KEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// Получение UserID из токена
func GetUserID(tokenString string) int64 {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(constants.SECRET_KEY), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token is valid")
	return claims.UserID
}
