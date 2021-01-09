package handlers

import (
	"github.com/EtienneBerube/only-cats/internal/models"
	"github.com/EtienneBerube/only-cats/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	req := models.LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAuth, err := services.GetUserAuthByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	token, err := services.ValidateLoginRequest(&req, userAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp := models.LoginResponse{
		ID:    userAuth.UserID,
		Token: token,
	}

	c.JSON(http.StatusOK, resp)
}

func SignUp(c *gin.Context) {
	req := models.SignUpRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.ValidateSignUpRequest(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		ID:            "",
		Name:          req.Name,
		Email:         req.Email,
		Subscriptions: []string{},
		Photos:        []string{},
		SubscriptionPrice: req.SubscriptionPrice,
		Balance: 100,
	}

	userId, err := services.CreateNewUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userAuth := models.UserAuth{
		Email:        user.Name,
		UserID:       userId,
		PasswordHash: services.GetPasswordHash(req.Email, req.Password),
	}

	_, err = services.CreateNewUserAuth(&userAuth)
	if err != nil {
		services.DeleteUser(userId) // Revert Creation
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := services.CreateToken(&userAuth)

	resp := models.SignUpResponse{
		ID:    userId,
		Token: token,
	}

	c.JSON(http.StatusCreated, resp)
}
