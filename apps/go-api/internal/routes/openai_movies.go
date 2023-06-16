package routes

import (
	"fmt"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
)

// GetMoviesFromGPT3RouteHandler is a gin route handler that takes a list of movie ids,
// fetches the movies, converts them to summaries, and sends those summaries to the GPT-3 API.
// @Summary Get movie suggestions based on favorite movies
// @Description This API receives a list of favorite movie IDs, fetches the corresponding movie summaries,
// and uses GPT-3 to generate movie suggestions based on these preferences.
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param movies body []int true "A list of favorite movie IDs"
// @Success 200 {array} types.Movie "Successful retrieval of movie suggestions"
// @Failure 400 {object} types.Error "Invalid input"
// @Failure 500 {object} types.Error "Failed to get movie summaries or suggestions"
// @Router /openai/movies [post]
func SuggestMoviesFromGPT3RouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonInput []int
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid input",
			})
			return
		}

		summaries, err := handlers.CreateMovieSummariesFromTMDBMovies(app, jsonInput)
		if err != nil {
			c.JSON(500, gin.H{
				// give more details in console
				"error": fmt.Sprintf("Failed to get movie summaries: %v", err),
			})
			return
		}
		movies, err := handlers.SuggestMoviesFromGPT3(app, summaries)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to get movie suggestions",
			})
			return
		}

		c.JSON(200, movies)
	}
}

func TranslateMoviesFromGPT3RouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the lang query parameter
		lang := c.Query("lang")

		// Accept a list of movie IDS, retrieve the Movies from TMDB, and translate the full Movie objects to the specified language
		var jsonInput []int
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid input",
			})
			return
		}

		// Initialize a slice to hold the movies
		var movies []types.Movie

		// Iterate over the movie IDs and retrieve each movie
		for _, id := range jsonInput {
			movie, err := handlers.FindMovieByID(app, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error": fmt.Sprintf("Failed to get movie with ID %d", id),
				})
				return
			}

			// Append the retrieved movie to the movies slice
			movies = append(movies, movie)
		}

		// we translate the movies to the specified language
		translatedMovies, err := handlers.TranslateMoviesFromGPT3(app, movies, lang)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("Failed to translate movies: %v", err),
			})
			return
		}

		c.JSON(200, translatedMovies)
	}
}
