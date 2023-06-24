package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	// Filter the response to only include the requested provider

	filteredAuthMethods := types.AuthMethodsResponse{}

	for _, authMethod := range authMethods.AuthProviders {
		if authMethod.Name == provider {
			authMethod.AuthURL = strings.Replace(authMethod.AuthURL, "redirect_uri=", "redirect_uri="+url.QueryEscape(app.OAuthRedirectURL), 1)
			filteredAuthMethods.AuthProviders = append(filteredAuthMethods.AuthProviders, authMethod)
		}
	}

	fmt.Println("Filtered auth methods:", filteredAuthMethods)

	return filteredAuthMethods, nil
}

// RefreshAuthToken refreshes an existing token
func RefreshAuthToken(app *types.App, token string) (types.RefreshResponse, error) {
	refreshResponse := types.RefreshResponse{}
	req, err := http.NewRequest("POST", app.PBAuthRefreshURL, nil)
	if err != nil {
		fmt.Println("RefreshAuthToken: NewRequest error:", err)
		return types.RefreshResponse{}, err
	}

	// Set the token in the request header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	err = utils.SendJSONRequest(req, &refreshResponse)
	if err != nil {
		fmt.Println("RefreshAuthToken: SendJSONRequest error:", err)
		return types.RefreshResponse{}, err
	}

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
