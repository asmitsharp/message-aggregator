package services

import (
	"github.com/matrix-org/gomatrix"
)

type MatrixService struct {
	HomeserverURL string
}

func NewMatrixService(homeserverURL string) *MatrixService {
	return &MatrixService{HomeserverURL: homeserverURL}
}

func (s *MatrixService) CreateAccount(username, password string) (*gomatrix.RespRegister, error) {
	client, err := gomatrix.NewClient(s.HomeserverURL, "", "")
	if err != nil {
		return nil, err
	}

	resp, err := client.RegisterDummy(&gomatrix.ReqRegister{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *MatrixService) Login(userID, accessToken string) (*gomatrix.Client, error) {
	client, err := gomatrix.NewClient(s.HomeserverURL, userID, accessToken)
	if err != nil {
		return nil, err
	}
	return client, nil
}
