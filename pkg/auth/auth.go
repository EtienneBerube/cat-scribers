package auth

import (
	"errors"
	"fmt"
	"github.com/EtienneBerube/cat-scribers/pkg/config"
	"github.com/dgrijalva/jwt-go"
)

// CreateToken creates a JWT token based on the claims provided by the caller.
func CreateToken(claims jwt.MapClaims) (string, error) {
	var err error

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(config.Config.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken validates a provided JWT token
func ValidateToken(token string) (bool, error) {
	result, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token: %s", token.Header["alg"])

		}
		return []byte(config.Config.JWTSecret), nil
	})

	if err != nil {
		return false, err
	}

	return result.Valid, nil
}

// ExtractUserId returns the authenticated user's ID from the JWT token provided
func ExtractUserId(token string) (string, error) {
	result, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token: %s", token.Header["alg"])

		}
		return []byte(config.Config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := result.Claims.(jwt.MapClaims); ok {
		id, _ := claims["user_id"].(string)
		return id, nil
	} else {
		return "", errors.New("Could not parse JWT claims")
	}
}
