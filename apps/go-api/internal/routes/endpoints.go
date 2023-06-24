// routes/endpoints.go
package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// WhoAmIRouteHandler returns information about the currently authenticated user.
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
func WhoAmIRouteHandler(app *types.App) gin.HandlerFunc {
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

// ListMoviesCollectionRouteHandler returns a list of movies from a specified collection.
// @Summary Get list of movies from movie collection
// @Description This API retrieves a list of movies from movie collection
// @Tags Movies
// @Accept  json
// @Produce  json
// @Security HTTPOnlySessionCookie
// @Success 200 {array} types.Movie "Successfully fetched movie collection"
// @Failure 500 {object} types.Error "Internal server error - Failed to get token or movies collection"
// @Router /movies [get]
func ListMoviesCollectionRouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log as admin to pb and get the token
		adminAuthResponse, err := services.PBAdminAuth(app)
		if err != nil {
			fmt.Println("ListMoviesCollection: Failed to get token", err)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get token",
			})
			return
		}

		token := adminAuthResponse.Token
		collectionUrl := app.MoviesCollectionURL

		var collection types.CollectionResponse
		err = services.PBGetCollection(collectionUrl, token, &collection)
		if err != nil {
			fmt.Println("ListMoviesCollection: Failed to get movies collection", err)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get movies collection",
			})
			return
		}

		var movies []types.Movie
		for _, item := range collection.Items {
			movieData, ok := item.(map[string]interface{})
			if !ok {
				fmt.Println("Failed to convert")
				continue
			}

			// Convert the map to a Movie
			jsonData, err := json.Marshal(movieData)
			if err != nil {
				fmt.Println("Failed to marshal movieData to json")
				continue
			}

			var movie types.Movie
			err = json.Unmarshal(jsonData, &movie)
			if err != nil {
				fmt.Println("Failed to unmarshal jsonData to movie")
				continue
			}

			movies = append(movies, movie)
		}

		c.JSON(http.StatusOK, movies) // send the movies to the client
	}
}

// PostUserMoviePreferencesRouteHandler creates a new user's movie preference record.
// @Summary Create new user's movie preferences
// @Description This API creates a new record in the user_has_movies collection
// @Tags User
// @Accept  json
// @Produce  json
// @Security HTTPOnlySessionCookie
// @Param userHasMovies body types.UserHasMovies true "User's movie preferences"
// @Success 200 {object} map[string]string "Successfully created user's movie preferences"
// @Failure 400 {object} types.Error "Bad request - No userId or token in session, Failed to parse request body"
// @Failure 500 {object} types.Error "Internal server error - Failed to update user's movie preferences"
// @Router /user/movie/preferences [post]
func PostUserMoviePreferencesRouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user's ID and token from the session
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

		// Get the user's movie preferences from the request body
		var userHasMovies types.UserHasMovies
		err := c.BindJSON(&userHasMovies)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to parse request body",
			})
			return
		}

		// Update the user's movie preferences in PocketBase
		err = handlers.UpdateUserHasMovies(app, userId, token, userHasMovies)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update user's movie preferences",
			})
			return
		}

		c.JSON(http.StatusOK, map[string]string{
			"message": "Successfully updated user's movie preferences",
		})
	}
}

// GetUserMoviePreferencesRouteHandler gets the user's movie preferences (see types.UserHasMovies)
// @Summary Get user's movie preferences
// @Description This API retrieves a user's movie preferences from user_has_movies collection
// @Tags User
// @Accept  json
// @Produce  json
// @Security HTTPOnlySessionCookie
// @Success 200 {object} types.UserHasMovies "Successfully fetched user's movie preferences"
// @Failure 400 {object} types.Error "Bad request - No userId or token in session"
// @Failure 500 {object} types.Error "Internal server error - Failed to get user's movie preferences"
// @Router /user/movie/preferences [get]
func GetUserMoviePreferencesRouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user's ID and token from the session
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

		// Get the user's movie preferences from PocketBase
		userHasMovies, err := handlers.GetUserHasMovies(app, userId, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get user's movie preferences",
			})
			return
		}

		c.JSON(http.StatusOK, userHasMovies)

	}
}
