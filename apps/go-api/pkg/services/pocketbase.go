package services

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"
)

type AdminAuthRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AdminAuthResponse struct {
	Token string `json:"token"`
	Admin struct {
		ID      string `json:"id"`
		Created string `json:"created"`
		Updated string `json:"updated"`
		Email   string `json:"email"`
		Avatar  int    `json:"avatar"`
	} `json:"admin"`
}

// PBAdminAuth authenticates as admin
func PBAdminAuth(app *types.App) (AdminAuthResponse, error) {
	adminAuthRequest := AdminAuthRequest{
		Identity: app.PBAuthAdmin,
		Password: app.PBAuthAdminPassword,
	}

	baseURL, err := url.Parse(app.PBUrl)
	if err != nil {
		return AdminAuthResponse{}, err
	}

	pathURL, err := url.Parse("/api/admins/auth-with-password")
	if err != nil {
		return AdminAuthResponse{}, err
	}

	adminAuthUrl := baseURL.ResolveReference(pathURL).String()

	adminAuthResponse := AdminAuthResponse{}

	// Send the request and decode the response
	err = utils.PostJSON(adminAuthUrl, adminAuthRequest, &adminAuthResponse) // Use = instead of :=
	if err != nil {
		return AdminAuthResponse{}, err
	}
	fmt.Println("adminAuthResponse:", adminAuthResponse)
	return adminAuthResponse, nil
}

// PBGetCollection retrieves a collection (including items) from PocketBase
func PBGetCollection(collectionUrl string, token string, out interface{}) error {
	req, err := http.NewRequest("GET", collectionUrl, nil)
	if err != nil {
		fmt.Println("PBGetCollection: Failed to create request", err)
		return err
	}

	// Add Authorization header with token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	err = utils.GetJSON(req, out)
	if err != nil {
		fmt.Println("PBGetCollection: GetJSON failed", err)
		return err
	}

	return nil
}
