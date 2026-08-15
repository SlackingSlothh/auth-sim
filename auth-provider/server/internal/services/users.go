package services

import (
	"github.com/gin-gonic/gin"
)

func RouteUsersAPI(router *gin.Engine) {
	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUser)
	router.POST("/users", registerUser)
	router.PATCH("/users/:id", updateUser)
	router.POST("/users/:id/groups", enlistUserToGroups)
	router.DELETE("/users/:id/groups", delistUserFromGroups)
}

// GET /users
func getAllUsers(c *gin.Context) {
}

// GET /users/:id
func getUser(c *gin.Context) {
}

// POST /users
func registerUser(c *gin.Context) {
}

// PATCH /users/:id
func updateUser(c *gin.Context) {
}

// POST /users/:id/groups
func enlistUserToGroups(c *gin.Context) {
}

// DELETE /users/:id/groups
func delistUserFromGroups(c *gin.Context) {
}