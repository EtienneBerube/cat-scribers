package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/repositories"
	"github.com/EtienneBerube/cat-scribers/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"regexp"
	"time"
)

const EMAIL_REGEX = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

// GetAuthByEmail returns the UserAuth associated with the email address
func GetUserAuthByEmail(email string) (*models.UserAuth, error) {
	return repositories.GetAuthByEmail(email)
}

// CreateNewUserAuth creates a new UserAuth
func CreateNewUserAuth(userAuth *models.UserAuth) (string, error) {
	id, err := repositories.SaveAuth(userAuth)
	if err != nil {
		return "", err
	}
	return id, nil
}

// DeleteUserAuth removes a UserAuth associated with the ID
func DeleteUserAuth(id string) error {
	return repositories.DeleteAuth(id)
}

// ModifyUserAuth modifies a UserAuth
func ModifyUserAuth(id string, userAuth *models.UserAuth) (bool, error) {
	return repositories.UpdateAuth(id, userAuth)
}

// CreateToken creates a JWT token for a given user
func CreateToken(userAuth *models.UserAuth) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userAuth.UserID
	claims["exp"] = time.Now().Add(time.Hour * 8760).Unix() // One year

	token, err := auth.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateSignUpRequest checks that the sign up request is valid
func ValidateSignUpRequest(req models.SignUpRequest) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	used, err := repositories.IsEmailUsed(req.Email)
	if err != nil {
		return err
	}

	if used {
		return errors.New("Email is already used by another user")
	}

	return nil
}

// ValidateLoginRequest checks that the login request is valid and that the credentials match
func ValidateLoginRequest(loginRequest *models.LoginRequest, userAuth *models.UserAuth) (token string, err error) {
	tryHash := GetPasswordHash(loginRequest.Email, loginRequest.Password)
	if (userAuth.Email == loginRequest.Email) && (tryHash == userAuth.PasswordHash) {
		token, err = CreateToken(userAuth)

		if err != nil {
			return "", err
		}

	} else {
		err = errors.New("Failed Login Attempt")
	}

	return token, err
}
// GetPasswordHash returns the hash for a given password. Uses email las prepend salt
func GetPasswordHash(email string, password string) string {
	msg := email + password
	hash := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", hash)
}

// validateEmail validates that the email matches the EMAIL_REGEX provided above. It is a modification of the RFC 5322 standard
func validateEmail(email string) error {
	matched, err := regexp.MatchString(EMAIL_REGEX, email)
	if err != nil {
		return err
	} else if !matched {
		return errors.New("This is not a valid email address")
	} else {
		return nil
	}
}
/* validatePassword Checks that the password matches certain conditions:
	- 6 characters or more
	- one or more digits
	- one or more uppercase letters
 */
func validatePassword(password string) error {
	var err error

	if len(password) < 6 {
		err = errors.New("Password must contain 6 or more characters")
	} else if matched, errs := regexp.MatchString("\\d+", password); errs != nil || !matched {
		err = errors.New("Password must contain at least one digit")
	} else if matched, errs := regexp.MatchString("[A-Z]", password); errs != nil || !matched {
		err = errors.New("Password must contain at least one capital letter")
	} else {
		err = nil
	}

	return err
}
