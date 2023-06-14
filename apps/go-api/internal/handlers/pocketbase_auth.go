package handlers

import (
	"fmt"
	"net/http"

	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"
)

// GetAuthMethods retrieves available authentication methods for a given provider
func GetAuthMethods(app *types.App, provider string) (types.AuthMethodsResponse, error) {
	authMethods := types.AuthMethodsResponse{}

	// Create a new request using http
	req, err := http.NewRequest("GET", app.PBAuthMethodsURL, nil)
	if err != nil {
		return types.AuthMethodsResponse{}, err
	}

	// Send the request and decode the response
	err = utils.GetJSON(req, &authMethods)
	if err != nil {
		return types.AuthMethodsResponse{}, err
	}

	//fmt print the response
	fmt.Println("GetAuthMethods: authMethods:", authMethods)
	// Filter the response to only include the requested provider

	filteredAuthMethods := types.AuthMethodsResponse{}

	for _, authMethod := range authMethods.AuthProviders {
		if authMethod.Name == provider {
			filteredAuthMethods.AuthProviders = append(filteredAuthMethods.AuthProviders, authMethod)
		}
	}

	return filteredAuthMethods, nil
}

// RefreshAuthToken refreshes an existing token
func RefreshAuthToken(app *types.App, token string) (types.RefreshResponse, error) {
	refreshResponse := types.RefreshResponse{}

	fmt.Println("RefreshAuthToken: PBAuthRefreshURL:", app.PBAuthRefreshURL)
	fmt.Println("RefreshAuthToken: token:", token)

	err := utils.PostJSON(app.PBAuthRefreshURL, token, &refreshResponse)
	if err != nil {
		fmt.Println("RefreshAuthToken: PostJSON error:", err)
		return types.RefreshResponse{}, err
	}

	fmt.Println("RefreshAuthToken: refreshResponse:", refreshResponse)
	return refreshResponse, nil
}

// TradeCodeForToken exchanges an authorization code for a token
func TradeCodeForToken(app *types.App, OAuthRequest types.OAuthRequest) (types.TradeResponse, error) {
	tradeResponse := types.TradeResponse{}
	err := utils.PostJSON(app.PBTradeURL, OAuthRequest, &tradeResponse)
	if err != nil {
		return types.TradeResponse{}, err
	}

	return tradeResponse, nil
}

// GetUserFromPocketBase retrieves user information for a given record id
func GetUserFromPocketBase(app *types.App, userId string, token string) (types.PocketBaseUserRecord, error) {
	userResponse := types.PocketBaseUserRecord{}

	// Create a new request
	req, err := http.NewRequest("GET", app.PBUserURL+userId, nil)
	if err != nil {
		return types.PocketBaseUserRecord{}, err
	}

	// Set the token in the request header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send the request and decode the response
	err = utils.GetJSON(req, &userResponse)
	if err != nil {
		return types.PocketBaseUserRecord{}, err
	}

	return userResponse, nil
}
