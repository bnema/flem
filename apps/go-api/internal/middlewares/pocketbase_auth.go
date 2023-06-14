package middlewares

import (
	"fmt"
	"net/http"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// IsLoggedIn is a middleware that checks if a user is logged in
func IsLoggedIn(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := app.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Session retrieval failed",
			})
			c.Abort()
			return
		}

		if session.Values["token"] == nil {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Not logged in",
			})
			c.Abort()
			return
		}
		c.Set("session", session)
		c.Next()
	}
}

// VerifyToken is a middleware that verifies the validity of a token
func VerifyToken(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := c.MustGet("session").(*sessions.Session)
		fmt.Printf("Session in VerifyToken: %+v\n", session)
		token := session.Values["token"].(string)
		fmt.Println("Token:", token)
		refreshResponse, err := handlers.RefreshAuthToken(app, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		session.Values["token"] = refreshResponse.Token // Save refreshed token in session
		err = session.Save(c.Request, c.Writer)         // Save session data
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to save session",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
