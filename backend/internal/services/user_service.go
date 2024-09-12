package services

import (
	"errors"

	"github.com/ashmitsharp/msg-agg/internal/models"
	"github.com/matrix-org/gomatrix"
	"gorm.io/gorm"
)

type UserService struct {
	DB                  *gorm.DB
	PlatformAuthService *PlatformAuthService
}

func NewUserService(db *gorm.DB, platformAuthService *PlatformAuthService) *UserService {
	return &UserService{DB: db, PlatformAuthService: platformAuthService}
}

func (s *UserService) CreateUser(email, password, name, matrixPassword string, matrixResp *gomatrix.RespRegister) (*models.User, error) {
	user := &models.User{
		Email:             email,
		Name:              name,
		MatrixUserID:      matrixResp.UserID,
		MatrixAccessToken: matrixResp.AccessToken,
		MatrixPassword:    matrixPassword,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	masterUser := &models.MasterUser{
		Name:  name,
		Email: email,
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(masterUser).Error; err != nil {
			return err
		}

		user.MasterUserID = masterUser.ID
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

func (s *UserService) AddIdentity(masterUserID, platform, identifier, credential string) (*models.UserIdentity, error) {
	verified, err := s.PlatformAuthService.VerifyIdentity(platform, identifier, credential)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, errors.New("failed to verify identity")
	}

	identity := &models.UserIdentity{
		MasterUserID: masterUserID,
		Platform:     platform,
		Identifier:   identifier,
		DisplayName:  identifier, // You might want to fetch the actual display name from the platform
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(identity).Error; err != nil {
			return err
		}

		// Fetch the MasterUser and append the new identity
		var masterUser models.MasterUser
		if err := tx.First(&masterUser, masterUserID).Error; err != nil {
			return err
		}

		masterUser.Identities = append(masterUser.Identities, *identity)
		return tx.Save(&masterUser).Error
	})

	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *UserService) GetIdentifier(masterUserID string, platform string) (*models.UserIdentity, error) {
	var slackIdentity models.UserIdentity
	if err := s.DB.Where("master_user_id = ? AND platform = ?", masterUserID, platform).First(&slackIdentity).Error; err != nil {
		return nil, err
	}
	return &slackIdentity, nil
}

func (s *UserService) GetIdentities(masterUserID string) ([]models.UserIdentity, error) {
	var identities []models.UserIdentity
	if err := s.DB.Where("master_user_id = ?", masterUserID).Find(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}
