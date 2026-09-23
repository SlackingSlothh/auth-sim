package services

import (
	"errors"
	"net/http"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/database"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/transport"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func RouteGroupsAPI(router *gin.RouterGroup) {
	router.GET("/groups", getAllGroups)
	router.GET("/groups/:id", getGroup)
	router.POST("/groups", createGroup)
	router.PATCH("/groups/:id", editGroup)
	router.POST("/groups/:id/members", addUsersToGroup)
	router.DELETE("/groups/:id/members", removeUsersFromGroup)
	router.GET("/groups/:id/nonmembers", getNonmembers)
	router.GET("/groups/:id/denied-apps", getDeniedApps)
	router.DELETE("/groups/:id/allowed-apps", removeAllowedApps)
	router.POST("/groups/:id/allowed-apps", addAllowedApps)
}

// GET /groups
func getAllGroups(c *gin.Context) {
	groups, err := database.SelectAllGroups()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
		return
	}

	AllGroupResponses := make([]transport.AllGroupResponse, len(groups))
	for i, g := range groups {
		AllGroupResponses[i].ID = g.ID.String()
		AllGroupResponses[i].Name = g.Name
		AllGroupResponses[i].Description = g.Description
	}

	c.JSON(http.StatusOK, AllGroupResponses)
}

// GET /groups/:id
func getGroup(c *gin.Context) {
	id := c.Param("id")
	var uuid datatypes.UUID
	if err := uuid.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrUserNotFound.Error()})
		return
	}

	group, apps, err := database.SelectGroup(uuid)
	if err == nil {
		var response transport.GetGroupResponse
		response.FromData(group, apps)
		c.JSON(http.StatusOK, gin.H{"group": response.Info, "users": response.Members, "apps": response.AllowedApps})
		return
	}
	
	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// POST /groups
func createGroup(c *gin.Context) {
	var newGroup models.Group
	if err := c.ShouldBindJSON(&newGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.CreateGroup(newGroup)
	var info transport.CreateGroupResponse
	info.ID = newGroup.ID.String()
	info.Name = newGroup.Name
	info.Description = newGroup.Description
	info.CreatedAt = newGroup.CreatedAt
	info.UpdatedAt = newGroup.UpdatedAt

	if err == nil {
		c.JSON(http.StatusCreated, info)
		return
	}

	if errors.Is(err, models.ErrGroupAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

}

// PATCH /groups/:id
func editGroup(c *gin.Context) {
	id := c.Param("id")
	var updatedInfo transport.UpdateGroupRequest

	if err := c.ShouldBindJSON(&updatedInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedGroup models.Group
	if err := updatedGroup.ID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}
	updatedGroup.Name = updatedInfo.Name
	updatedGroup.Description = updatedInfo.Description

	err := database.UpdateGroupData(updatedGroup)

	if err == nil {
		c.JSON(http.StatusOK, updatedInfo)
		return
	}

	if errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

}

// GET /groups/:id/nonmembers
func getNonmembers(c *gin.Context) {
	id := c.Param("id")

	var groupID datatypes.UUID
	if err := groupID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}

	users, err := database.GetGroupNonmembers(groupID)
	if err == nil {
		c.JSON(http.StatusOK, users)
		return
	}

	if errors.Is(err, models.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

}

// GET /groups/:id/denied-apps
func getDeniedApps(c *gin.Context) {
	id := c.Param("id")

	var groupID datatypes.UUID
	if err := groupID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}

	apps, err := database.GetDeniedApps(groupID)
	if err == nil {
		c.JSON(http.StatusOK, apps)
		return
	}

	if errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// DELETE /groups/:id/allowed-apps
func removeAllowedApps(c *gin.Context) {
	id := c.Param("id")
	var groupID datatypes.UUID
	
	if err := groupID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}

	var appIDs []string
	if err := c.ShouldBindJSON(&appIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var appUUIDs []datatypes.UUID
	for _, gid := range appIDs {
		var theUuid datatypes.UUID
		if err := theUuid.Scan(gid); err == nil {
			appUUIDs = append(appUUIDs, theUuid)
		}
	}

	removedApps, err := database.RemoveAllowedApps(groupID, appUUIDs)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"removed": removedApps})
		return
	}

	if errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

}

// POST /groups/:id/allowed-apps
func addAllowedApps(c *gin.Context) {
	id := c.Param("id")
	var groupID datatypes.UUID
	
	if err := groupID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}

	var appIDs []string
	if err := c.ShouldBindJSON(&appIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var appUUIDs []datatypes.UUID
	for _, gid := range appIDs {
		var theUuid datatypes.UUID
		if err := theUuid.Scan(gid); err == nil {
			appUUIDs = append(appUUIDs, theUuid)
		}
	}

	addedApps, err := database.AppendAllowedApps(groupID, appUUIDs)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"added": addedApps})
		return
	}

	if errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// POST /groups/:id/members
func addUsersToGroup(c *gin.Context) {
	id := c.Param("id")
	var groupID datatypes.UUID
	
	if err := groupID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}

	var userIDs []string
	if err := c.ShouldBindJSON(&userIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userUUIDs []datatypes.UUID
	for _, gid := range userIDs {
		var theUuid datatypes.UUID
		if err := theUuid.Scan(gid); err == nil {
			userUUIDs = append(userUUIDs, theUuid)
		}
	}

	addedGroups, err := database.AddGroupMembers(groupID, userUUIDs)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"added": addedGroups})
		return
	}

	if errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// DELETE /groups/:id/members
func removeUsersFromGroup(c *gin.Context) {
	id := c.Param("id")
	var groupID datatypes.UUID
	
	if err := groupID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrGroupNotFound.Error()})
		return
	}

	var userIDs []string
	if err := c.ShouldBindJSON(&userIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userUUIDs []datatypes.UUID
	for _, gid := range userIDs {
		var theUuid datatypes.UUID
		if err := theUuid.Scan(gid); err == nil {
			userUUIDs = append(userUUIDs, theUuid)
		}
	}

	removedUsers, err := database.RemoveGroupMembers(groupID, userUUIDs)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"removed": removedUsers})
		return
	}

	if errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}