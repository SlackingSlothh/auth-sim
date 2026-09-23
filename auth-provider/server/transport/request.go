package transport

type RegisterUserRequest struct {
	Name     string
	Email    string
	Password string
}

type UpdateGroupRequest struct {
	Name        string
	Description string
}

type RegisterAppRequest struct {
	Name                  string  `json:"name"`
	LogoutNotificationURL string  `json:"logoutNotificationUrl"`
	LaunchURL             *string `json:"launchUrl"`
}

type UpdateAppRequest struct {
	Name                  string  `json:"name"`
	LogoutNotificationURL string  `json:"logout_notification_url"`
	LaunchURL             *string `json:"launch_url"`
	Status                string  `json:"status"`
}

type RedirectURIRequest struct {
	RedirectURI string `json:"redirectUri"`
}

type EditRedirectURIRequest struct {
	OldRedirectURI string `json:"oldRedirectUri"`
	NewRedirectURI string `json:"newRedirectUri"`
}

type LoginRequest struct {
	Email    string
	Password string
}

type ChangePasswordRequest struct {
	Password string
}