// routes/endpoints.go
package routes

import (
	"net/http"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

// WhoAmI returns information about the currently authenticated user.
// It retrieves the user's ID and token from the session,
// and uses these to fetch the user's data from PocketBase.
// @Summary Get current user information
// @Description This API retrieves information about the currently authenticated user
// by fetching the user's data from PocketBase using the user's ID and token found in the session.
// @Tags User
// @Accept  json
// @Produce  json
// @Security HTTPOnlySessionCookie
// @Success 200 {object} types.PocketBaseUserRecord "Successfully fetched user data"
// @Failure 400 {object} types.Error "Invalid request - No userId or token in session"
// @Failure 500 {object} types.Error "Internal server error - Failed to get user from PocketBase"
// @Router /whoami [get]
func WhoAmI(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := c.MustGet("session").(*sessions.Session)

		userId, ok := session.Values["userId"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": "No userId in session",
			})
			return
		}

		token, ok := session.Values["token"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": "No token in session",
			})
			return
		}

		user, err := handlers.GetUserFromPocketBase(app, userId, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get user from PocketBase",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
