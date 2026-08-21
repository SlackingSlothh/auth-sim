package transport

import (
	"time"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
)

type GroupRecord struct {
	ID          string
	Name        string
	Description string
}

type GroupInfo struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (this *GroupInfo) InfoOf(group models.Group) {
	this.ID = group.ID.String()
	this.Name = group.Name
	this.Description = group.Description
	this.CreatedAt = group.CreatedAt
	this.UpdatedAt = group.UpdatedAt
}

type Nonmember struct {
	ID    string
	Name  string
	Email string
}

type DeniedApps struct {
	ID   string
	Name string
}

type AppRecord struct {
	ID       string
	Name     string
	ClientID string
	Status   string
}

type AppInfo struct {
	Name                  string  `json:"name"`
	ClientID              string  `json:"client_id"`
	Status                string  `json:"status"`
	LaunchURL             *string `json:"launch_url"`
	LogoutNotificationURL string  `json:"logout_notification_url"`
}

type AppAllowedGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AppDetail struct {
	App           AppInfo           `json:"app"`
	AllowedGroups []AppAllowedGroup `json:"allowedGroups"`
	RedirectURIs  []string          `json:"redirectURIs"`
}

type RegisterAppRequest struct {
	Name                  string  `json:"name"`
	LogoutNotificationURL string  `json:"logoutNotificationUrl"`
	LaunchURL             *string `json:"launchUrl"`
}

type RegisterAppResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RedirectURIRequest struct {
	RedirectURI string `json:"redirectUri"`
}

type EditRedirectURIRequest struct {
	OldRedirectURI string `json:"oldRedirectUri"`
	NewRedirectURI string `json:"newRedirectUri"`
}
