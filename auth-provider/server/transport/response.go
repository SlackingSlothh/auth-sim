package transport

import (
	"time"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
)

type AllGroupResponse struct {
	ID          string
	Name        string
	Description string
}

type groupInfo struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type groupMember struct {
	ID    string
	Name  string
	Email string
}

type allowedApp struct {
	ID     string
	Name   string
	Status string
}

type GetGroupResponse struct {
	Info		groupInfo
	Members 	[]groupMember
	AllowedApps []allowedApp
}

func (resp *GetGroupResponse) FromData(group models.Group, apps []models.Application) {
	resp.Info = groupInfo{
		ID: group.ID.String(),
		Name: group.Name,
		Description: group.Description,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}
	resp.Members = make([]groupMember, len(group.Users))
	for i, g := range group.Users {
		resp.Members[i].ID = g.ID.String()
		resp.Members[i].Name = g.Name
		resp.Members[i].Email = g.Email
	}
	resp.AllowedApps = make([]allowedApp, len(apps))
	for i, g := range apps {
		resp.AllowedApps[i].ID = g.ID.String()
		resp.AllowedApps[i].Name = g.Name
		resp.AllowedApps[i].Status = g.Status
	}
}

type CreateGroupResponse struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RegisterAppResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
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

type AppAllowedGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AppDetail struct {
	App AppInfo `json:"app"`
}

type AppInfo struct {
	ID                    string            `json:"id"`
	Name                  string            `json:"name"`
	ClientID              string            `json:"client_id"`
	Status                string            `json:"status"`
	LaunchURL             *string           `json:"launch_url"`
	LogoutNotificationURL string            `json:"logout_notification_url"`
	RedirectURIs          []string          `json:"redirect_uris"`
	AllowedGroups         []AppAllowedGroup `json:"allowed_groups"`
}