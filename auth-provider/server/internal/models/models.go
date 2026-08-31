package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string         `gorm:"type:varchar(255);not null"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Status       string         `gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;autoUpdateTime;not null"`

	Groups []Group `gorm:"many2many:user_groups;joinForeignKey:UserID;joinReferences:GroupID"`
}

func (User) TableName() string {
	return "users"
}

type Group struct {
	ID          datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;autoUpdateTime;not null"`

	Users []User `gorm:"many2many:user_groups;joinForeignKey:GroupID;joinReferences:UserID"`
}

func (Group) TableName() string {
	return "groups"
}

type UserGroup struct {
	ID        datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    datatypes.UUID `gorm:"type:uuid;not null;index:idx_user_group_unique,unique"`
	GroupID   datatypes.UUID `gorm:"type:uuid;not null;index:idx_user_group_unique,unique"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`

	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Group Group `gorm:"foreignKey:GroupID;references:ID"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}

type Application struct {
	ID                    datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name                  string         `gorm:"type:varchar(255);not null"`
	ClientID              string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	ClientSecretHash      *string        `gorm:"type:varchar(255)"`
	Status                string         `gorm:"type:varchar(50);not null"`
	LaunchURL             *string        `gorm:"type:text"`
	LogoutNotificationURL string         `gorm:"type:text;not null"`
	CreatedAt             time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	UpdatedAt             time.Time      `gorm:"type:timestamptz;autoUpdateTime;not null"`
}

func (Application) TableName() string {
	return "applications"
}

type ApplicationRedirectURI struct {
	ID            datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ApplicationID string         `gorm:"type:varchar(255);not null;index"`
	RedirectURI   string         `gorm:"type:text;not null"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`

	Application Application `gorm:"foreignKey:ApplicationID;references:ClientID"`
}

func (ApplicationRedirectURI) TableName() string {
	return "application_redirect_uris"
}

type ApplicationGroupPolicy struct {
	ID            datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ApplicationID datatypes.UUID `gorm:"type:uuid;not null;index:idx_app_group_policy_unique,unique"`
	GroupID       datatypes.UUID `gorm:"type:uuid;not null;index:idx_app_group_policy_unique,unique"`
	Effect        string         `gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`

	Application Application `gorm:"foreignKey:ApplicationID;references:ID"`
	Group       Group       `gorm:"foreignKey:GroupID;references:ID"`
}

func (ApplicationGroupPolicy) TableName() string {
	return "application_group_policies"
}

type SSOSession struct {
	ID              datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          datatypes.UUID `gorm:"type:uuid;not null;index"`
	SessionTokenHash string       `gorm:"type:varchar(255);not null"`
	Status          string         `gorm:"type:varchar(50);not null"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	ExpiresAt       time.Time      `gorm:"type:timestamptz;not null"`
	LastActivityAt  *time.Time     `gorm:"type:timestamptz"`
	RevokedAt       *time.Time     `gorm:"type:timestamptz"`
	RevokeReason    *string        `gorm:"type:varchar(255)"`
	IPAddress       *string        `gorm:"type:varchar(45)"`
	UserAgent       *string        `gorm:"type:text"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (SSOSession) TableName() string {
	return "sso_sessions"
}

type AuthorizationCode struct {
	ID            datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CodeHash      string         `gorm:"type:varchar(255);not null"`
	UserID        datatypes.UUID `gorm:"type:uuid;not null;index"`
	ApplicationID string         `gorm:"type:varchar(255);not null;index"`
	SSOSessionID  datatypes.UUID `gorm:"type:uuid;not null;index"`
	RedirectURI   string         `gorm:"type:text;not null"`
	Status        string         `gorm:"type:varchar(50);not null"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	ExpiresAt     time.Time      `gorm:"type:timestamptz;not null"`
	UsedAt        *time.Time     `gorm:"type:timestamptz"`

	User        User        `gorm:"foreignKey:UserID;references:ID"`
	Application Application `gorm:"foreignKey:ApplicationID;references:ClientID"`
	SSOSession  SSOSession  `gorm:"foreignKey:SSOSessionID;references:ID"`
}

func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}

type AccessToken struct {
	ID            datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TokenHash     string         `gorm:"type:varchar(255);not null"`
	UserID        datatypes.UUID `gorm:"type:uuid;not null;index"`
	ApplicationID string         `gorm:"type:varchar(255);not null;index"`
	SSOSessionID  datatypes.UUID `gorm:"type:uuid;not null;index"`
	Scopes        datatypes.JSON `gorm:"type:jsonb"`
	Status        string         `gorm:"type:varchar(50);not null"`
	IssuedAt      time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`
	ExpiresAt     time.Time      `gorm:"type:timestamptz;not null"`
	RevokedAt     *time.Time     `gorm:"type:timestamptz"`

	User        User        `gorm:"foreignKey:UserID;references:ID"`
	Application Application `gorm:"foreignKey:ApplicationID;references:ClientID"`
	SSOSession  SSOSession  `gorm:"foreignKey:SSOSessionID;references:ID"`
}

func (AccessToken) TableName() string {
	return "access_tokens"
}

type AuditLog struct {
	ID            datatypes.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventType     string         `gorm:"type:varchar(255);not null"`
	ActorID       *datatypes.UUID `gorm:"type:uuid"`
	UserID        *datatypes.UUID `gorm:"type:uuid;index"`
	ApplicationID *string        `gorm:"type:varchar(255);index"`
	SessionID     *datatypes.UUID `gorm:"type:uuid;index"`
	Result        string         `gorm:"type:varchar(50);not null"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`
	IPAddress     *string        `gorm:"type:varchar(45)"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;autoCreateTime;not null"`

	User        *User        `gorm:"foreignKey:UserID;references:ID"`
	Application *Application `gorm:"foreignKey:ApplicationID;references:ClientID"`
	Session     *SSOSession  `gorm:"foreignKey:SessionID;references:ID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// =======================================================================

type UserInfo struct {
	ID           string
	Name         string
	Email        string
	Status       string
}

func (this User) Info() UserInfo {
	return UserInfo{
		ID: this.ID.String(),
		Name: this.Name,
		Email: this.Email,
		Status: this.Status,
	}
}

func (this *User) FromInfo(info UserInfo) {
	this.ID.Scan(info.ID)
	if info.Email != "" {
		this.Email = info.Email
	}
	if info.Name != "" {
		this.Name = info.Name
	}
	if info.Status != "" {
		this.Status = info.Status
	}
}

type GroupBrief struct {
	ID string
	Name string
}