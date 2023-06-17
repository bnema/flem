package routes

import (
	"fmt"
	"net/http"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
)

// LoginRoute godoc
// @Summary Initiate OAuth authentication
// @Description This route handles the '/login' endpoint and initiates OAuth authentication
// @Tags OAuth
// @Accept  json
// @Produce  json
// @Param provider query string true "OAuth provider name"
// @Success 302 {string} string "Redirection to the OAuth URL"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [get]
func LoginRoute(app *types.App, c *gin.Context) {
	provider := c.Query("provider")
	authMethods, err := handlers.GetAuthMethods(app, provider) // Get authentication methods for the specified provider
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Failed to get auth methods: %v", err),
		})
		return
	}

	session, err := app.SessionStore.Get(c.Request, "session") // Get session for this request
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get session: %v", err),
		})
		return
	}
	session.Options.MaxAge = 60 * 15 // Set session max age to 15 minutes
	session.Options.HttpOnly = true  // Set session cookie to HTTP only
	session.Values["provider"] = provider
	session.Values["state"] = authMethods.AuthProviders[0].State
	session.Values["codeVerifier"] = authMethods.AuthProviders[0].CodeVerifier
	session.Values["authUrl"] = authMethods.AuthProviders[0].AuthURL
	session.Save(c.Request, c.Writer) // Save session data
	fmt.Println("Saved provider:", session.Values["provider"])
	fmt.Println("Saved state:", session.Values["state"])
	fmt.Println("Saved codeVerifier:", session.Values["codeVerifier"])
	fmt.Println("Saved authUrl:", session.Values["authUrl"])

	c.Redirect(302, authMethods.AuthProviders[0].AuthURL) // Redirect user to OAuth URL
}

// RedirectRoute godoc
// @Summary Finalize OAuth authentication
// @Description This route handles the '/oauth-redirect' endpoint and finalizes the OAuth authentication process. After successful authentication, the session is updated with a token and a userId.
// @Tags OAuth
// @Accept  json
// @Produce  html
// @Param code query string true "OAuth code received from provider"
// @Param state query string true "OAuth state received from provider"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /oauth-redirect [get]
func RedirectRoute(app *types.App, c *gin.Context) {
	session, _ := app.SessionStore.Get(c.Request, "session") // Get session for this request
	provider := session.Values["provider"].(string)
	state := session.Values["state"].(string)
	codeVerifier := session.Values["codeVerifier"].(string)
	code := c.Query("code")

	if state != c.Query("state") {
		c.JSON(400, map[string]string{
			"error": "Invalid state",
		})
		return
	}

	oAuthRequest := types.OAuthRequest{
		Provider:     provider,
		Code:         code,
		CodeVerifier: codeVerifier,
		RedirectURL:  app.OAuthRedirectURL,
		State:        state,
	}

	tradeResponse, err := handlers.TradeCodeForToken(app, oAuthRequest)
	if err != nil {
		c.JSON(400, map[string]string{
			"error": "Failed to trade code for token",
		})
		return
	}

	session.Values["token"] = tradeResponse.Token
	session.Values["userId"] = tradeResponse.Record.Id

	// Save session data
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(500, map[string]string{
			"error": "Failed to save session",
		})
		return
	}

	// close the page
	c.Data(200, "text/html", []byte("You can close this page now"))
}
