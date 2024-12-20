package routes

import (
	"fmt"
	"net/http"
	"sync"

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
// @Security HTTPOnlySessionCookie
// @Router /openai/movies [post]
func SuggestMoviesFromGPT3RouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonInput []int
		var maxInputLength int = 10
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid input",
			})
			return
		}

		// Limit the number of movies to 10
		if len(jsonInput) > maxInputLength {
			// Error JSON
			c.JSON(400, gin.H{
				"error": fmt.Sprintf("Too many movies. Please limit your input to %d movies", maxInputLength),
			})
			return

		}

		summaries, err := handlers.CreateMovieSummariesFromTMDBMovies(app, jsonInput)
		fmt.Println(summaries)
		if err != nil {
			c.JSON(500, gin.H{
				// give more details in console
				"error": fmt.Sprintf("Failed to get movie summaries: %v", err),
			})
			return
		}
		movies, err := handlers.SuggestMoviesFromGPT3(app, summaries)
		fmt.Println(movies)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Failed to get movie suggestions",
			})
			return
		}

		c.JSON(200, movies)
	}
}

// TranslateMoviesFromGPT3RouteHandler is a gin route handler that takes a list of movie ids,
// fetches the movies, and translates them to the specified language using GPT-3.
// @Summary Translate movies to a specified language
// @Description This API receives a list of movie IDs and translates the corresponding movie information to the specified language.
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param lang query string true "The language to translate to"
// @Param movies body []int true "A list of movie IDs"
// @Success 200 {array} types.Movie "Successful translation of movies"
// @Failure 400 {object} types.Error "Invalid input"
// @Failure 500 {object} types.Error "Failed to get movie with ID"
// @Security HTTPOnlySessionCookie
// @Router /openai/translate [post]
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
		var mu sync.Mutex
		var wg sync.WaitGroup
		errCh := make(chan error, 1) // Create error channel

		// Iterate over the movie IDs and retrieve each movie
		for _, id := range jsonInput {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				movie, err := handlers.FindMovieByID(app, id)
				if err != nil {
					errCh <- fmt.Errorf("failed to get movie with ID %d: %v", id, err)
					return
				}
				// Append the retrieved movie to the movies slice
				mu.Lock()
				movies = append(movies, movie)
				mu.Unlock()
			}(id)
		}

		wg.Wait()

		select {
		case err := <-errCh: // Check if there's an error
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("Failed to process movie list: %v", err),
			})
			return
		default: // If there's no error, continue
		}

		translatedMovies, existingMovies, err := handlers.TranslateMoviesFromGPT3(app, movies, lang)
		if err != nil {
			c.JSON(500, gin.H{
				"error": fmt.Sprintf("Failed to translate movies: %v", err),
			})
			return
		}

		// Add existing movies to the response if they exist
		if len(existingMovies) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"existing_movies":   existingMovies,
				"translated_movies": translatedMovies,
			})
			return
		}

		// If there are no existing movies, just return the translated movies
		c.JSON(http.StatusOK, translatedMovies)
	}
}
