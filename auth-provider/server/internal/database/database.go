package database

import (
	"errors"

	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/config"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/models"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/transport"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SelectAllUsers() ([]models.User, error) {
	db := config.GetDB()
	var users []models.User

	result := db.Find(&users)

	if result.Error != nil {
		return users, models.ErrInternal
	}

	return users, nil
}

func SelectUser(id datatypes.UUID) (models.User, error) {
	db := config.GetDB()
	var user models.User

	result := db.Preload("Groups").First(&user, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, models.ErrUserNotFound
		} else {
			return user, models.ErrInternal
		}
	}

	return user, nil
}

func CreateUser(newUser models.User) error {
	db := config.GetDB()
	var existingUser models.User

	result := db.Where(models.User{Email: newUser.Email}).
		Omit("ID").
		Attrs(newUser).
		FirstOrCreate(&existingUser)

	if result.RowsAffected == 0 {
		return models.ErrUserAlreadyExists
	}
	if result.Error != nil {
		return models.ErrInternal
	}
	return nil
}

func UpdateUserData(user models.User) error {
	db := config.GetDB()

	updates := map[string]interface{}{}

	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.Status != "" {
		updates["status"] = user.Status
	}
	if user.PasswordHash != "" {
		updates["password_hash"] = user.PasswordHash
	}

	result := db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(updates)

	if result.RowsAffected == 0 {
		return models.ErrUserNotFound
	}
	if result.Error != nil {
		return models.ErrInternal
	}
	return nil
}

func AppendUserGroups(userId datatypes.UUID, groupIds []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(groupIds) == 0 {
		return []datatypes.UUID{}, nil
	}

	records := make([]models.UserGroup, 0, len(groupIds))
	for _, groupID := range groupIds {
		records = append(records, models.UserGroup{
			UserID:  userId,
			GroupID: groupID,
		})
	}

	err := db.Table("user_groups").
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "group_id"}},
				DoNothing: true,
			},
		).
		Create(&records).Error

	if err == nil {
		return groupIds, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrUserNotFound
	}

	return nil, models.ErrInternal
}

func RemoveUserGroups(userId datatypes.UUID, groupIds []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(groupIds) == 0 {
		return []datatypes.UUID{}, nil
	}

	err := db.Model(&models.UserGroup{}).
		Where("user_id = ? AND group_id IN ?", userId, groupIds).
		Delete(&models.UserGroup{}).Error

	if err == nil {
		return groupIds, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrUserNotFound
	}

	return nil, models.ErrInternal
}

func SelectUserByEmail(email string) (models.User, error) {
	db := config.GetDB()
	var user models.User

	result := db.First(&user, "email = ?", email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, models.ErrUserNotFound
		} else {
			return user, models.ErrInternal
		}
	}

	return user, nil
}

// func GetUserFromSession(cookie string) (models.User, error) {
// 	db := config.GetDB()
// 	var session models.SSOSession
// 	session.
// }

func IsAllowedOrigin(origin string) bool {
	db := config.GetDB()
	var app models.ApplicationRedirectURI
	err := db.Where("redirect_uri LIKE ?", origin+"%").First(&app)
	return err == nil
}

func GetAvailableGroups(userID datatypes.UUID) ([]models.GroupBrief, error) {
	db := config.GetDB()

	var user models.User
	result := db.Select("id").First(&user, "id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrInternal
	}

	var groups []models.GroupBrief
	result = db.Table("groups").
		Select("groups.id AS id, groups.name AS name").
		Joins("LEFT JOIN user_groups ON user_groups.group_id = groups.id AND user_groups.user_id = ?", userID).
		Where("user_groups.group_id IS NULL").
		Order("groups.name ASC").
		Scan(&groups)

	if result.Error != nil {
		return nil, models.ErrInternal
	}

	return groups, nil
}

func SelectAllGroups() ([]models.Group, error) {
	db := config.GetDB()
	var groups []models.Group

	result := db.Find(&groups)

	if result.Error != nil {
		return nil, models.ErrInternal
	}

	return groups, nil
}

func SelectAllApps() ([]transport.AppRecord, error) {
	db := config.GetDB()
	var apps []transport.AppRecord

	result := db.Model(&models.Application{}).
		Select("id, name, client_id, status").
		Find(&apps)

	if result.Error != nil {
		return nil, models.ErrInternal
	}

	return apps, nil
}

func SelectApp(id datatypes.UUID) (transport.AppDetail, error) {
	db := config.GetDB()
	var app models.Application
	var detail transport.AppDetail

	result := db.First(&app, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return detail, models.ErrAppNotFound
		}
		return detail, models.ErrInternal
	}

	detail.App = transport.AppInfo{
		Name:                  app.Name,
		ClientID:              app.ClientID,
		Status:                app.Status,
		LaunchURL:             app.LaunchURL,
		LogoutNotificationURL: app.LogoutNotificationURL,
	}

	result = db.Table("groups").
		Select("groups.id AS id, groups.name AS name").
		Joins("JOIN application_group_policies ON application_group_policies.group_id = groups.id").
		Where("application_group_policies.application_id = ?", app.ClientID).
		Order("groups.name ASC").
		Scan(&detail.AllowedGroups)
	if result.Error != nil {
		return detail, models.ErrInternal
	}

	result = db.Model(&models.ApplicationRedirectURI{}).
		Where("application_id = ?", app.ClientID).
		Order("redirect_uri ASC").
		Pluck("redirect_uri", &detail.RedirectURIs)
	if result.Error != nil {
		return detail, models.ErrInternal
	}

	return detail, nil
}

func CreateApp(newApp models.Application) error {
	db := config.GetDB()
	var existingApp models.Application

	result := db.Where(models.Application{Name: newApp.Name}).
		Omit("ID").
		Attrs(newApp).
		FirstOrCreate(&existingApp)

	if result.RowsAffected == 0 {
		return models.ErrAppAlreadyExists
	}
	if result.Error != nil {
		return models.ErrInternal
	}

	return nil
}

func AppendAppGroups(appID datatypes.UUID, groupIDs []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(groupIDs) == 0 {
		return []datatypes.UUID{}, nil
	}

	var app models.Application
	result := db.Select("client_id").First(&app, "id = ?", appID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrAppNotFound
		}
		return nil, models.ErrInternal
	}

	records := make([]models.ApplicationGroupPolicy, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		records = append(records, models.ApplicationGroupPolicy{
			ApplicationID: app.ClientID,
			GroupID:       groupID,
			Effect:        "allow",
		})
	}

	err := db.Table("application_group_policies").
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "application_id"}, {Name: "group_id"}, {Name: "effect"}},
				DoNothing: true,
			},
		).
		Create(&records).Error

	if err == nil {
		return groupIDs, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrGroupNotFound
	}

	return nil, models.ErrInternal
}

func RemoveAppGroups(appID datatypes.UUID, groupIDs []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(groupIDs) == 0 {
		return []datatypes.UUID{}, nil
	}

	var app models.Application
	result := db.Select("client_id").First(&app, "id = ?", appID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrAppNotFound
		}
		return nil, models.ErrInternal
	}

	err := db.Model(&models.ApplicationGroupPolicy{}).
		Where("application_id = ? AND group_id IN ? AND effect = ?", app.ClientID, groupIDs, "allow").
		Delete(&models.ApplicationGroupPolicy{}).Error

	if err != nil {
		return nil, models.ErrInternal
	}

	return groupIDs, nil
}

func CreateAppRedirectURI(appID datatypes.UUID, redirectURI string) (string, error) {
	db := config.GetDB()

	clientID, err := selectAppClientID(db, appID)
	if err != nil {
		return "", err
	}

	var existingURI models.ApplicationRedirectURI
	result := db.Where(models.ApplicationRedirectURI{
		ApplicationID: clientID,
		RedirectURI:   redirectURI,
	}).FirstOrCreate(&existingURI)

	if result.Error != nil {
		return "", models.ErrInternal
	}
	if result.RowsAffected == 0 {
		return "", models.ErrRedirectURIAlreadyExists
	}

	return redirectURI, nil
}

func UpdateAppRedirectURI(appID datatypes.UUID, oldRedirectURI string, newRedirectURI string) (string, error) {
	db := config.GetDB()

	clientID, err := selectAppClientID(db, appID)
	if err != nil {
		return "", err
	}

	if oldRedirectURI != newRedirectURI {
		var existingURI models.ApplicationRedirectURI
		result := db.Select("id").
			First(&existingURI, "application_id = ? AND redirect_uri = ?", clientID, newRedirectURI)
		if result.Error == nil {
			return "", models.ErrRedirectURIAlreadyExists
		}
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", models.ErrInternal
		}
	}

	result := db.Model(&models.ApplicationRedirectURI{}).
		Where("application_id = ? AND redirect_uri = ?", clientID, oldRedirectURI).
		Update("redirect_uri", newRedirectURI)

	if result.Error != nil {
		return "", models.ErrInternal
	}
	if result.RowsAffected == 0 {
		return "", models.ErrRedirectURINotFound
	}

	return newRedirectURI, nil
}

func DeleteAppRedirectURI(appID datatypes.UUID, redirectURI string) (string, error) {
	db := config.GetDB()

	clientID, err := selectAppClientID(db, appID)
	if err != nil {
		return "", err
	}

	result := db.Where("application_id = ? AND redirect_uri = ?", clientID, redirectURI).
		Delete(&models.ApplicationRedirectURI{})

	if result.Error != nil {
		return "", models.ErrInternal
	}
	if result.RowsAffected == 0 {
		return "", models.ErrRedirectURINotFound
	}

	return redirectURI, nil
}

func selectAppClientID(db *gorm.DB, appID datatypes.UUID) (string, error) {
	var app models.Application
	result := db.Select("client_id").First(&app, "id = ?", appID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", models.ErrAppNotFound
		}
		return "", models.ErrInternal
	}

	return app.ClientID, nil
}

func CreateGroup(newGroup models.Group) error {
	db := config.GetDB()
	var existingGroup models.Group

	result := db.Where(models.Group{Name: newGroup.Name}).
		Omit("ID").
		Attrs(newGroup).
		FirstOrCreate(&existingGroup)

	if result.RowsAffected == 0 {
		return models.ErrGroupAlreadyExists
	}
	if result.Error != nil {
		return models.ErrInternal
	}
	return nil
}

func SelectGroup(id datatypes.UUID) (models.Group, []models.Application, error) {
	db := config.GetDB()
	var group models.Group

	result := db.Preload("Users").First(&group, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return group, nil, models.ErrGroupNotFound
		} else {
			return group, nil, models.ErrInternal
		}
	}

	var policies []models.ApplicationGroupPolicy
	err := db.
		Preload("Application").
		Where("group_id = ?", group.ID).
		Find(&policies).Error

	if err != nil {
		return group, nil, models.ErrInternal
	}

	apps := make([]models.Application, 0, len(policies))
	for _, policy := range policies {
		apps = append(apps, policy.Application)
	}

	return group, apps, nil
}

func UpdateGroupData(group models.Group) error {
	db := config.GetDB()

	updates := map[string]interface{}{}

	if group.Name != "" {
		updates["name"] = group.Name
	}
	if group.Description != "" {
		updates["description"] = group.Description
	}

	result := db.Model(&models.Group{}).
		Where("id = ?", group.ID).
		Updates(updates)

	if result.RowsAffected == 0 {
		return models.ErrGroupNotFound
	}
	if result.Error != nil {
		return models.ErrInternal
	}
	return nil
}

func GetGroupNonmembers(groupID datatypes.UUID) ([]transport.Nonmember, error) {
	db := config.GetDB()

	var group models.Group
	result := db.Select("id").First(&group, "id = ?", groupID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrGroupNotFound
		}
		return nil, models.ErrInternal
	}

	var users []transport.Nonmember
	result = db.Table("users").
		Select("users.id AS id, users.name AS name, email").
		Joins("LEFT JOIN user_groups ON user_groups.group_id = ? AND user_groups.user_id = users.id", groupID).
		Where("user_groups.user_id IS NULL").
		Order("users.name ASC").
		Scan(&users)

	if result.Error != nil {
		return nil, models.ErrInternal
	}

	return users, nil
}

func GetDeniedApps(groupID datatypes.UUID) ([]transport.DeniedApps, error) {
	db := config.GetDB()

	var group models.Group
	result := db.Select("id").First(&group, "id = ?", groupID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrGroupNotFound
		}
		return nil, models.ErrInternal
	}

	var apps []transport.DeniedApps
	result = db.Table("applications").
		Select("applications.id AS id, applications.name AS name").
		Joins("LEFT JOIN application_group_policies ON application_group_policies.applications_id = applications.id AND application_group_policies.group_id = ?", groupID).
		Where("application_group_policies.group_id IS NULL").
		Order("applications.name ASC").
		Scan(&apps)

	if result.Error != nil {
		return nil, models.ErrInternal
	}

	return apps, nil
}

func RemoveAllowedApps(groupID datatypes.UUID, appIDs []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(appIDs) == 0 {
		return []datatypes.UUID{}, nil
	}

	err := db.Model(&models.ApplicationGroupPolicy{}).
		Where("application_id IN ? AND group_id = ?", appIDs, groupID).
		Delete(&models.ApplicationGroupPolicy{}).Error

	if err == nil {
		return appIDs, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrGroupNotFound
	}

	return nil, models.ErrInternal
}

func RemoveGroupMembers(groupID datatypes.UUID, userIDs []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(userIDs) == 0 {
		return []datatypes.UUID{}, nil
	}

	err := db.Model(&models.UserGroup{}).
		Where("user_id IN ? AND group_id = ?", userIDs, groupID).
		Delete(&models.UserGroup{}).Error

	if err == nil {
		return userIDs, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrGroupNotFound
	}

	return nil, models.ErrInternal
}

func AppendAllowedApps(groupID datatypes.UUID, appIDs []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(appIDs) == 0 {
		return []datatypes.UUID{}, nil
	}

	records := make([]models.ApplicationGroupPolicy, 0, len(appIDs))
	for _, appID := range appIDs {
		records = append(records, models.ApplicationGroupPolicy{
			ApplicationID: appID.String(),
			GroupID:       groupID,
		})
	}

	err := db.Table("application_group_policy").
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "application_id"}, {Name: "group_id"}},
				DoNothing: true,
			},
		).
		Create(&records).Error

	if err == nil {
		return appIDs, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrGroupNotFound
	}

	return nil, models.ErrInternal
}

func AddGroupMembers(groupID datatypes.UUID, userIDs []datatypes.UUID) ([]datatypes.UUID, error) {
	db := config.GetDB()

	if len(userIDs) == 0 {
		return []datatypes.UUID{}, nil
	}

	records := make([]models.UserGroup, 0, len(userIDs))
	for _, userID := range userIDs {
		records = append(records, models.UserGroup{
			UserID:  userID,
			GroupID: groupID,
		})
	}

	err := db.Table("user_groups").
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "group_id"}},
				DoNothing: true,
			},
		).
		Create(&records).Error

	if err == nil {
		return userIDs, nil
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return nil, models.ErrGroupNotFound
	}

	return nil, models.ErrInternal
}
