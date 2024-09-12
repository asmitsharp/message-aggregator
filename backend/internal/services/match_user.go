package services

import (
	"strings"

	"github.com/ashmitsharp/msg-agg/internal/models"
	"github.com/ashmitsharp/msg-agg/internal/utils"
	"gorm.io/gorm"
)

func MatchOrCreateUser(db *gorm.DB, platform, identifier, displayName string) (*models.MasterUser, error) {
	var identity models.UserIdentity
	if err := db.Where("platform = ? AND identifier = ?", platform, identifier).First(&identity).Error; err == nil {
		// Existing identity found, return associated MasterUser
		var masterUser models.MasterUser
		if err := db.Preload("Identities").First(&masterUser, identity.MasterUserID).Error; err != nil {
			return nil, err
		}
		return &masterUser, nil
	}

	// No existing identity, try to match based on email or phone
	var masterUser models.MasterUser
	if strings.Contains(identifier, "@") {
		// Treat identifier as email
		if err := db.Where("email = ?", identifier).First(&masterUser).Error; err == nil {
			// Matching MasterUser found, create new identity
			identity = models.UserIdentity{
				MasterUserID: masterUser.ID,
				Platform:     platform,
				Identifier:   identifier,
				DisplayName:  displayName,
			}
			if err := db.Create(&identity).Error; err != nil {
				return nil, err
			}
			masterUser.Identities = append(masterUser.Identities, identity)
			return &masterUser, nil
		}
	}

	// No match found, create new MasterUser and UserIdentity
	var email, phone string
	if strings.Contains(identifier, "@") {
		email = identifier
		phone = ""
	} else {
		email = ""
		phone = identifier
	}

	masterUser = models.MasterUser{
		ID:    utils.GenerateUUID(),
		Name:  displayName,
		Email: email,
		Phone: phone,
	}
	if err := db.Create(&masterUser).Error; err != nil {
		return nil, err
	}

	identity = models.UserIdentity{
		MasterUserID: masterUser.ID,
		Platform:     platform,
		Identifier:   identifier,
		DisplayName:  displayName,
	}
	if err := db.Create(&identity).Error; err != nil {
		return nil, err
	}

	masterUser.Identities = []models.UserIdentity{identity}
	return &masterUser, nil
}
