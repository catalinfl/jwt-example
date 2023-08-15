package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type RequestUser struct {
	Name     string
	Password string
}

func CreateJWTToken(user *RequestUser) (string, int64, error) {
	exp := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}

func VerifyToken(c *fiber.Ctx) error {
	authHeader := c.Cookies("jwt")

	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing auth token")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid auth token")
	}

	if !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	c.Locals("userID", claims["name"])

	return c.Next()
}
