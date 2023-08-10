package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Create a new JWT Token
func CreateToken(userID uint, role string) (string, error) {
	// Get .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create Claims
	claims := jwt.MapClaims{
		"authorized": true,
		"role":       role,
		"userID":     userID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

// Create a new JWT Cookie Token
func CreateCookieToken(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:    "fiberStore",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
		Path:    "/",
	})
}

// Delete JWT Cookie Token
func DeleteCookieToken(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    "fiberStore",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
}
