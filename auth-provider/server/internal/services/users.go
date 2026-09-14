package services

import (
	"errors"
	"net/http"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/config"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/database"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/transport"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func RouteUsersAPI(router *gin.RouterGroup) {
	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", registerUser)
	router.PATCH("/users/:id", updateUser)
	router.POST("/users/:id/groups", enlistUserToGroups)
	router.DELETE("/users/:id/groups", delistUserFromGroups)
	router.PUT("/users/:id/password", changePassword)
	router.GET("/users/:id/available-groups", getAvailableGroups)
}

// GET /users
func getAllUsers(c *gin.Context) {
	users, err := database.SelectAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
		return
	}

	usersInfo := make([]models.UserInfo, len(users))
	for i, u := range users {
		usersInfo[i].ID = u.ID.String()
		usersInfo[i].Email = u.Email
		usersInfo[i].Name = u.Name
		usersInfo[i].Status = u.Status
	}

	c.JSON(http.StatusOK, usersInfo)
}

// GET /users/:id
func getUser(c *gin.Context) {
	id := c.Param("id")
	var uuid datatypes.UUID
	if err := uuid.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrUserNotFound.Error()})
		return
	}

	user, err := database.SelectUser(uuid)
	if err == nil {
		groupBriefs := make([]models.GroupBrief, len(user.Groups))
		for i, g := range user.Groups {
			groupBriefs[i].ID = g.ID.String()
			groupBriefs[i].Name = g.Name
		}
		c.JSON(http.StatusOK, gin.H{"user": user.Info(), "groups": groupBriefs})
		return
	}
	
	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// GET /users/:id/available-groups
// func getUserAvailableGroups(c *gin.Context) {
// 	id := c.Param("id")
// }

// POST /users
func registerUser(c *gin.Context) {
	var request transport.RegisterUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := config.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUser := models.User{
		Name: request.Name,
		Email: request.Email,
		PasswordHash: passwordHash,
		Status: "active",
	}
	
	err = database.CreateUser(newUser)

	if err == nil {
		c.JSON(http.StatusCreated, newUser.Info())
		return
	}

	if errors.Is(err, models.ErrUserAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// PATCH /users/:id
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedInfo models.UserInfo

	updatedInfo.ID = id
	if err := c.ShouldBindJSON(&updatedInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedUser models.User
	updatedUser.FromInfo(updatedInfo)
	err := database.UpdateUserData(updatedUser)

	if err == nil {
		c.JSON(http.StatusOK, updatedInfo)
		return
	}

	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

}

// POST /users/:id/groups
func enlistUserToGroups(c *gin.Context) {
	id := c.Param("id")
	var userId datatypes.UUID
	
	if err := userId.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrUserNotFound.Error()})
		return
	}

	var groupIds []string
	if err := c.ShouldBindJSON(&groupIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var groupUuids []datatypes.UUID
	for _, gid := range groupIds {
		var theUuid datatypes.UUID
		if err := theUuid.Scan(gid); err == nil {
			groupUuids = append(groupUuids, theUuid)
		}
	}

	addedGroups, err := database.AppendUserGroups(userId, groupUuids)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"added": addedGroups})
		return
	}

	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// DELETE /users/:id/groups
func delistUserFromGroups(c *gin.Context) {
	id := c.Param("id")
	var userId datatypes.UUID
	
	if err := userId.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrUserNotFound.Error()})
		return
	}

	var groupIds []string
	if err := c.ShouldBindJSON(&groupIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var groupUuids []datatypes.UUID
	for _, gid := range groupIds {
		var theUuid datatypes.UUID
		if err := theUuid.Scan(gid); err == nil {
			groupUuids = append(groupUuids, theUuid)
		}
	}

	removedGroups, err := database.RemoveUserGroups(userId, groupUuids)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"removed": removedGroups})
		return
	}

	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// PUT /users/:id/password
func changePassword(c *gin.Context) {
	id := c.Param("id")
	var uuid datatypes.UUID
	if err := uuid.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrUserNotFound.Error()})
		return
	}

	type passwordRequest struct {
		Password string
	}

	var req passwordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.SelectUser(uuid)
	
	if err == nil {
		user.PasswordHash, err = config.HashPassword(req.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = database.UpdateUserData(user)
	}

	if err == nil {
		c.Status(http.StatusOK)
		return
	}

	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// TODO: REVOKE CENTRAL SESSION
}

func getAvailableGroups(c *gin.Context) {
	id := c.Param("id")

	var userID datatypes.UUID
	if err := userID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrUserNotFound.Error()})
		return
	}

	groups, err := database.GetAvailableGroups(userID)
	if err == nil {
		c.JSON(http.StatusOK, groups)
		return
	}

	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}