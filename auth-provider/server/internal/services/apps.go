package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/config"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/database"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/transport"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func RouteAppsAPI(router *gin.RouterGroup) {
	router.GET("/apps", getAllApps)
	router.GET("/apps/:id", getApp)
	router.POST("/apps", registerApp)
	router.PATCH("/apps/:id", updateApp)
	router.GET("/apps/:id/available-groups", getAppAvailableGroups)
	router.POST("/apps/:id/groups", allowGroupsToApp)
	router.DELETE("/apps/:id/groups", denyGroupsFromApp)
	router.POST("/apps/:id/redirect-uris", addAppRedirectURI)
	router.PATCH("/apps/:id/redirect-uris", editAppRedirectURI)
	router.DELETE("/apps/:id/redirect-uris", deleteAppRedirectURI)
}

// GET /apps
func getAllApps(c *gin.Context) {
	apps, err := database.SelectAllApps()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
		return
	}

	c.JSON(http.StatusOK, apps)
}

// GET /apps/:id
func getApp(c *gin.Context) {
	id := c.Param("id")
	var appID datatypes.UUID
	if err := appID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrAppNotFound.Error()})
		return
	}

	app, err := database.SelectApp(appID)
	if err == nil {
		c.JSON(http.StatusOK, app)
		return
	}

	if errors.Is(err, models.ErrAppNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// POST /apps
func registerApp(c *gin.Context) {
	var req transport.RegisterAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.LogoutNotificationURL = strings.TrimSpace(req.LogoutNotificationURL)
	if req.LaunchURL != nil {
		launchURL := strings.TrimSpace(*req.LaunchURL)
		if launchURL == "" {
			req.LaunchURL = nil
		} else {
			req.LaunchURL = &launchURL
		}
	}

	if req.Name == "" || req.LogoutNotificationURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and logoutNotificationUrl are required"})
		return
	}

	clientID, err := config.GenerateRandomShit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrInternal.Error()})
		return
	}

	clientSecret, err := config.GenerateRandomShit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrInternal.Error()})
		return
	}

	clientSecretHash, err := config.HashPassword(clientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrInternal.Error()})
		return
	}

	newApp := models.Application{
		Name:                  req.Name,
		ClientID:              clientID,
		ClientSecretHash:      &clientSecretHash,
		Status:                "active",
		LaunchURL:             req.LaunchURL,
		LogoutNotificationURL: req.LogoutNotificationURL,
	}

	err = database.CreateApp(newApp)
	if err == nil {
		c.JSON(http.StatusCreated, transport.RegisterAppResponse{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		})
		return
	}

	if errors.Is(err, models.ErrAppAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// PATCH /apps/:id
func updateApp(c *gin.Context) {
	appID, ok := parseAppID(c)
	if !ok {
		return
	}

	var req transport.UpdateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.LogoutNotificationURL = strings.TrimSpace(req.LogoutNotificationURL)
	req.Status = strings.TrimSpace(req.Status)
	if req.LaunchURL != nil {
		launchURL := strings.TrimSpace(*req.LaunchURL)
		if launchURL == "" {
			req.LaunchURL = nil
		} else {
			req.LaunchURL = &launchURL
		}
	}

	if req.Name == "" || req.LogoutNotificationURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and logout_notification_url are required"})
		return
	}

	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be active or inactive"})
		return
	}

	err := database.UpdateAppData(appID, req.Name, req.LaunchURL, req.LogoutNotificationURL, req.Status)
	if err == nil {
		c.JSON(http.StatusOK, req)
		return
	}

	if errors.Is(err, models.ErrAppNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// GET /apps/:id/available-groups
func getAppAvailableGroups(c *gin.Context) {
	appID, ok := parseAppID(c)
	if !ok {
		return
	}

	groups, err := database.GetAppAvailableGroups(appID)
	if err == nil {
		c.JSON(http.StatusOK, groups)
		return
	}

	if errors.Is(err, models.ErrAppNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// POST /apps/:id/groups
func allowGroupsToApp(c *gin.Context) {
	id := c.Param("id")
	var appID datatypes.UUID
	if err := appID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrAppNotFound.Error()})
		return
	}

	var groupIDs []string
	if err := c.ShouldBindJSON(&groupIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupUUIDs, processedIDs := parseUUIDStrings(groupIDs)
	addedGroups, err := database.AppendAppGroups(appID, groupUUIDs)
	if err == nil {
		if len(addedGroups) == len(processedIDs) {
			c.JSON(http.StatusOK, gin.H{"added": processedIDs})
			return
		}

		c.JSON(http.StatusOK, gin.H{"added": uuidStrings(addedGroups)})
		return
	}

	if errors.Is(err, models.ErrAppNotFound) || errors.Is(err, models.ErrGroupNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// DELETE /apps/:id/groups
func denyGroupsFromApp(c *gin.Context) {
	id := c.Param("id")
	var appID datatypes.UUID
	if err := appID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrAppNotFound.Error()})
		return
	}

	var groupIDs []string
	if err := c.ShouldBindJSON(&groupIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupUUIDs, processedIDs := parseUUIDStrings(groupIDs)
	removedGroups, err := database.RemoveAppGroups(appID, groupUUIDs)
	if err == nil {
		if len(removedGroups) == len(processedIDs) {
			c.JSON(http.StatusOK, gin.H{"removed": processedIDs})
			return
		}

		c.JSON(http.StatusOK, gin.H{"removed": uuidStrings(removedGroups)})
		return
	}

	if errors.Is(err, models.ErrAppNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// POST /apps/:id/redirect-uris
func addAppRedirectURI(c *gin.Context) {
	appID, ok := parseAppID(c)
	if !ok {
		return
	}

	var req transport.RedirectURIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redirectURI := strings.TrimSpace(req.RedirectURI)
	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "redirectUri is required"})
		return
	}

	addedURI, err := database.CreateAppRedirectURI(appID, redirectURI)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{"added": addedURI})
		return
	}

	writeRedirectURIError(c, err)
}

// PATCH /apps/:id/redirect-uris
func editAppRedirectURI(c *gin.Context) {
	appID, ok := parseAppID(c)
	if !ok {
		return
	}

	var req transport.EditRedirectURIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldRedirectURI := strings.TrimSpace(req.OldRedirectURI)
	newRedirectURI := strings.TrimSpace(req.NewRedirectURI)
	if oldRedirectURI == "" || newRedirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "oldRedirectUri and newRedirectUri are required"})
		return
	}

	updatedURI, err := database.UpdateAppRedirectURI(appID, oldRedirectURI, newRedirectURI)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"updated": updatedURI})
		return
	}

	writeRedirectURIError(c, err)
}

// DELETE /apps/:id/redirect-uris
func deleteAppRedirectURI(c *gin.Context) {
	appID, ok := parseAppID(c)
	if !ok {
		return
	}

	var req transport.RedirectURIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redirectURI := strings.TrimSpace(req.RedirectURI)
	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "redirectUri is required"})
		return
	}

	deletedURI, err := database.DeleteAppRedirectURI(appID, redirectURI)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"deleted": deletedURI})
		return
	}

	writeRedirectURIError(c, err)
}

func parseUUIDStrings(ids []string) ([]datatypes.UUID, []string) {
	uuids := make([]datatypes.UUID, 0, len(ids))
	processedIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		var uuid datatypes.UUID
		if err := uuid.Scan(id); err == nil {
			uuids = append(uuids, uuid)
			processedIDs = append(processedIDs, id)
		}
	}

	return uuids, processedIDs
}

func uuidStrings(ids []datatypes.UUID) []string {
	strings := make([]string, len(ids))
	for i, id := range ids {
		strings[i] = id.String()
	}
	return strings
}

func parseAppID(c *gin.Context) (datatypes.UUID, bool) {
	id := c.Param("id")
	var appID datatypes.UUID
	if err := appID.Scan(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrAppNotFound.Error()})
		return appID, false
	}

	return appID, true
}

func writeRedirectURIError(c *gin.Context, err error) {
	if errors.Is(err, models.ErrAppNotFound) || errors.Is(err, models.ErrRedirectURINotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if errors.Is(err, models.ErrRedirectURIAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
