package routes

import (
	"math/rand"
	"net/url"
	"strconv"

	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TmdbApiResponse struct {
	Results []types.TmdbMovie `json:"results"`
}

// @Summary Search movies by title
// @Description Get movies that match given titles
// @Tags TMDB
// @ID get-movies-by-title
// @Accept json
// @Produce json
// @Param titles body []string true "List of Titles"
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /tmdb/movies/post/title [post]
func TMDBMovieByTitleRouteHandler(c *gin.Context) {
	titles := c.PostFormArray("titles")
	for _, title := range titles {
		var movies []types.Movie
		query := url.Values{}
		query.Add("query", title)
		err := services.CallTMDBApi("/search/movie", query, &movies)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Something went wrong",
			})
			return
		}
		c.JSON(200, movies)
	}
}

// @Summary Get movies by IDs
// @Description Get movies with given IDs
// @Tags TMDB
// @ID get-movies-by-ids
// @Accept json
// @Produce json
// @Param ids body []integer true "List of Movie IDs"
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /tmdb/movies/post/ids [post]
func TMDBMoviesByIDSRouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonInput []int
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid input",
			})
			return
		}

		var movies []types.Movie
		// Iterate over ids in the jsonInput
		for _, id := range jsonInput {
			movie, err := handlers.FindMovieByID(app, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error": "Something went wrong",
				})
				return
			}
			// Save each movie in the collection of movies
			handlers.SaveMovieToPocketbase(app, movie)
			movies = append(movies, movie)
		}

		c.JSON(200, movies)
	}
}

// @Summary Get random popular movies
// @Description Get 10 random popular movies
// @Tags TMDB
// @ID get-random-movies
// @Accept json
// @Produce json
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /tmdb/movies/random10 [get]
func TMDBRandomMoviesRouteHandler(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiResponse TmdbApiResponse
		query := url.Values{}

		// Generate a random page number between 1 and 500
		randomPage := strconv.Itoa(rand.Intn(500) + 1)
		query.Add("page", randomPage)

		// Random sort orders
		sortOrders := []string{"popularity.asc", "popularity.desc", "release_date.asc", "release_date.desc", "revenue.asc", "revenue.desc", "vote_average.asc", "vote_average.desc"}
		randomSortOrder := sortOrders[rand.Intn(len(sortOrders))]
		query.Add("sort_by", randomSortOrder)

		err := services.CallTMDBApi("/discover/movie", query, &apiResponse)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Something went wrong",
			})
			return
		}

		var movies []types.Movie
		for _, tmdbMovie := range apiResponse.Results {
			// Use the convert function to convert the tmdbMovie to a Movie
			movie := utils.ConvertTmdbMovieToMovie(tmdbMovie)

			// Save the movie to Pocketbase
			result, err := handlers.SaveMovieToPocketbase(app, movie)
			if err != nil {
				c.JSON(500, gin.H{
					"error": "Error saving movie to Pocketbase",
				})
				return
			}
			// Append the result of the save operation to the list of results
			movies = append(movies, result)
		}

		c.JSON(200, movies)
	}
}

// @Summary Get movies by genre and release date
// @Description Get movies that match the specified genre and were released in a specific year
// @Tags TMDB
// @ID get-movies-by-genre-date
// @Accept json
// @Produce json
// @Param genre query string true "Genre ID"
// @Param year query string false "Release Year"
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /tmdb/movies [get]
func TMDBMoviesByGenreAndDateRouteHandler(c *gin.Context) {
	genre := c.Query("genre")
	year := c.Query("year")
	var movies []types.Movie
	query := url.Values{}
	query.Add("with_genres", genre)
	query.Add("primary_release_year", year)
	result := &types.MovieDiscoveryResponse{}
	err := services.CallTMDBApi("/discover/movie", query, result)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	movies = result.Results
	c.JSON(200, movies)
}
