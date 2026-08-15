package services

import (
	"github.com/gin-gonic/gin"
)

func RouteGroupsAPI(router *gin.Engine) {
	router.GET("/groups", getAllGroups)
	router.GET("/groups/:id", getGroup)
	router.POST("/groups", createGroup)
	router.PATCH("/groups/:id", editGroup)
	router.POST("/groups/:id/members", addUsersToGroup)
	router.DELETE("/groups/:id/members", removeUsersFromGroup)
}

// GET /groups
func getAllGroups(c *gin.Context) {
}

// GET /groups/:id
func getGroup(c *gin.Context) {
}

// POST /groups
func createGroup(c *gin.Context) {
}

// PATCH /groups/:id
func editGroup(c *gin.Context) {
}

// POST /groups/:id/members
func addUsersToGroup(c *gin.Context) {
}

// DELETE /groups/:id/members
func removeUsersFromGroup(c *gin.Context) {
}