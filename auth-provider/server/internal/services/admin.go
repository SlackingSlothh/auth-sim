package services

import (
	"errors"
	"net/http"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/config"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/database"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
	"github.com/gin-gonic/gin"
)

type loginCredential struct {
	email string
	password string
}

// POST /login
func AdminLogin(c *gin.Context) {
	var credential loginCredential
	if err := c.ShouldBindJSON(&credential); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userFound, err := database.SelectUserByEmail(credential.email)
	
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if config.IsPasswordMatch(credential.password, userFound.PasswordHash) {
		c.JSON(http.StatusOK, gin.H{"user": userFound})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrInvalidCredential.Error()})
}

// func GetCurrentUser(c *gin.Context) {
// 	sessionCookie, err := c.Cookie("session")

// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "not authenticated",
// 		})
// 		return
// 	}

// 	// Validate sessionCookie against your database/session store.
// 	user, err := database.GetUserFromSession(sessionCookie)

// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "invalid session",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"id":    user.ID,
// 		"email": user.Email,
// 		"name":  user.Name,
// 	})
// }