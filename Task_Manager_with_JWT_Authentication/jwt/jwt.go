package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(userID int) (string, error) {
	if userID <= 0 {
		return "", fmt.Errorf("invalid userID")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKey := getEnv("SECRET_KEY", "mysecretkey")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (int, error) {
	if tokenString == "" {
		return 0, fmt.Errorf("tokenString is empty")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := getEnv("SECRET_KEY", "mysecretkey")
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := strconv.Atoi(fmt.Sprintf("%.0f", claims["userID"])) // Извлекаем userID из токена
		if err != nil {
			return 0, err
		}
		if userID <= 0 {
			return 0, fmt.Errorf("invalid userID in token")
		}
		return userID, nil
	}
	return 0, fmt.Errorf("invalid token")
}

// getEnv получает значение переменной окружения по имени или возвращает значение по умолчанию, если переменная не установлена.
func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}
