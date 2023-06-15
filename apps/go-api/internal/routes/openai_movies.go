package routes

import (
	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
)

// GetMoviesFromGPT3RouteHandler is a gin route handler that takes a list of movie ids,
// fetches the movies, converts them to summaries, and sends those summaries to the GPT-3 API
func GetMoviesFromGPT3RouteHandler(app *types.App) gin.HandlerFunc {
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
				"error": "Failed to get movie summaries",
			})
			return
		}
		movies, err := handlers.GetMoviesFromGPT3(app, summaries)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to get movie suggestions",
			})
			return
		}

		c.JSON(200, movies)
	}
}
