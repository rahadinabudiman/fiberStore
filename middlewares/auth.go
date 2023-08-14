package middlewares

import (
	"errors"
	"fiberStore/helpers"
	"log"
	"net/http"
	"os"
	"strings"
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

func GetTokenFromHeader(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		// If the header has "Bearer", split and get the token
		if strings.Contains(authHeader, "Bearer") {
			splitHeader := strings.Split(authHeader, " ")
			if len(splitHeader) == 2 {
				return splitHeader[1]
			}
		} else {
			// If it doesn't have "Bearer", return the whole header value as the token
			return authHeader
		}
	}
	return ""
}

func GetUserIdFromToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid claims")
	}

	getUserId, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("userID claim not found")
	}
	userID := uint(getUserId)
	return userID, nil
}

func GetRoleFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	getRole, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role claim not found")
	}

	return getRole, nil
}

func RoleMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		// Check if the user's role matches the required role
		if userRole != "Admin" && userRole != role {
			// Return an error response indicating unauthorized access
			errorResponse := helpers.ErrorResponses{
				StatusCode: fiber.StatusForbidden,
				Message:    "Forbidden",
				Errors:     "Unauthorized access",
			}
			return c.Status(http.StatusForbidden).JSON(errorResponse)
		}

		return c.Next()
	}
}

func JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.ErrorResponses{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Unauthorized",
			Errors:     "Missing Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil || !token.Valid {
		return JWTErrorHandler(err, c)
	}

	// Set the validated token in the context
	c.Locals("user", token)

	return c.Next()
}

func JWTErrorHandler(err error, c *fiber.Ctx) error {
	// Customize the JWT error response
	customError := helpers.ErrorResponses{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Errors:     err.Error(),
	}
	return c.Status(customError.StatusCode).JSON(customError)
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
