package services

import (
	"github.com/gin-gonic/gin"
)

func RouteAppsAPI(router *gin.Engine) {
	router.GET("/apps", getAllApps)
	router.GET("/apps/:id", getApp)
	router.POST("/apps", registerApp)
	router.POST("/apps/:id/groups", allowGroupsToApp)
	router.DELETE("/apps/:id/groups", denyGroupsFromApp)
}

// GET /apps
func getAllApps(c *gin.Context) {
}

// GET /apps/:id
func getApp(c *gin.Context) {
}

// POST /apps
func registerApp(c *gin.Context) {
}

// POST /apps/:id/groups
func allowGroupsToApp(c *gin.Context) {
}

// DELETE /apps/:id/groups
func denyGroupsFromApp(c *gin.Context) {
}