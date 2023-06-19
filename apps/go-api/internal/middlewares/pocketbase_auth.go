package middlewares

import (
	"net/http"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func SessionMiddleware(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := app.SessionStore.Get(c.Request, "session")
		if err != nil {
			// Return 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Session retrieval failed",
			})
			c.Abort()
			return

		}

		// store session data in gin.Context
		c.Set("session", session)

		c.Next()
	}
}

// VerifyToken is a middleware that verifies the validity of a token
func VerifyToken(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := c.MustGet("session").(*sessions.Session)
		tokenValue, ok := session.Values["token"]
		if !ok {
			// No token in session
			c.Set("error", "No token in session")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": c.MustGet("error")})
			return
		}
		token, ok := tokenValue.(string)
		if !ok {
			// Token in session is not a string
			c.Set("error", "Token in session is not a string")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": c.MustGet("error")})
			return
		}
		refreshResponse, err := handlers.RefreshAuthToken(app, token)
		if err != nil {
			c.Set("error", "Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": c.MustGet("error")})
			return
		}

		session.Values["token"] = refreshResponse.Token // Save refreshed token in session
		// We ensure that the session cookie is secure and HTTP only
		session.Options.Secure = true
		session.Options.HttpOnly = true

		err = session.Save(c.Request, c.Writer) // Save session data
		if err != nil {
			c.Set("error", "Failed to save session")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": c.MustGet("error")})
			return
		}
		c.Next()
	}
}
