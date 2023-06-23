package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// PBSaveItemToCollection saves an item to a collection in PocketBase
func PBSaveItemToCollection(collectionUrl string, token string, item interface{}) (interface{}, error) {
	var savedItem interface{}
	err := utils.PostJSON(collectionUrl, item, &savedItem, token)
	if err != nil {
		var dbError types.PBCollectionError
		// Try to unmarshal the error into a DBError
		if jsonErr := json.Unmarshal([]byte(err.Error()), &dbError); jsonErr == nil {
			errMessages := []string{}
			for key, value := range dbError.Data {
				message, ok := value["message"].(string) // Assert message to string
				if ok {
					errMessages = append(errMessages, fmt.Sprintf("%s: %s", key, message))
				}
			}
			fmt.Printf("PBSaveItemToCollection: Failed with code: %s. Errors: %s\n", dbError.Code, strings.Join(errMessages, ", "))
			return nil, nil // return nil error to continue the process if needed
		} else {
			fmt.Println("PBSaveItemToCollection: PostJSON failed", err)
			return nil, err
		}
	}
	return savedItem, nil
}

// PBGetItemFromCollection retrieves an item from a collection in PocketBase
func PBGetItemFromCollection(collectionUrl string, token string, filter string) (*types.CollectionResponse, error) {
	var resp types.CollectionResponse

	// Build the URL with the filters
	url := fmt.Sprintf("%s?filter=%s", collectionUrl, filter)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("PBGetItemFromCollection: Failed to create request", err)
		return nil, err
	}

	// Add Authorization header with token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	err = utils.GetJSON(req, &resp)
	if err != nil {
		fmt.Println("PBGetItemFromCollection: GetJSON failed", err)
		return nil, err
	}

	return &resp, nil
}

// // PBUpdateItemInCollection updates an item in a collection in PocketBase
// func PBUpdateItemInCollection(collectionUrl string, token string, item interface{}, out interface{}) error {
// 	err := utils.PutJSON(collectionUrl, item, out)
// 	if err != nil {
// 		fmt.Println("PBUpdateItemInCollection: PutJSON failed", err)
// 		return err
// 	}

// 	return nil
// }

// // PBDeleteItemFromCollection deletes an item from a collection in PocketBase
// func PBDeleteItemFromCollection(collectionUrl string, token string, out interface{}) error {
// 	err := utils.DeleteJSON(collectionUrl, out)
// 	if err != nil {
// 		fmt.Println("PBDeleteItemFromCollection: Delete failed", err)
// 		return err
// 	}

// 	return nil
// }
