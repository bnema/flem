package routes

import (
	"fmt"
	"net/url"

	"github.com/bnema/flem/go-api/internal/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/bnema/flem/go-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TmdbApiResponse struct {
	Results []types.Movie `json:"results"`
}

// @Summary Search movies by title
// @Description Get movies that match given titles
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
// @ID get-movies-by-ids
// @Accept json
// @Produce json
// @Param ids body []integer true "List of Movie IDs"
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /tmdb/movies/post/ids [post]
func TMDBMoviesByIDSRouteHandler(c *gin.Context) {
	var jsonInput []int
	if err := c.BindJSON(&jsonInput); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid input",
		})
		return
	}

	// Iterate over ids in the jsonInput
	for _, id := range jsonInput {
		var movie types.Movie
		// Declare a slice to hold summary movies
		var summaryMovies []types.SummaryItemMovie
		err := services.CallTMDBApi(fmt.Sprintf("/movie/%d", id), url.Values{}, &movie)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Something went wrong",
			})
			return
		} else if services.ValidateMovieData(movie) != nil {
			// we want another movie
			continue
		}
		// Create summary for each movie
		summaryMovie := utils.SummaryFromMovie(movie)
		summaryMovies = append(summaryMovies, summaryMovie)

		// Return summary movies as fmt
		fmt.Println(summaryMovies)

		c.JSON(200, movie)
	}
}

// @Summary Get random popular movies
// @Description Get 10 random popular movies
// @ID get-random-movies
// @Accept json
// @Produce json
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /tmdb/movies/random10 [get]
func TMDBRandomMoviesRouteHandler(c *gin.Context) {
	var apiResponse TmdbApiResponse
	query := url.Values{}
	query.Add("sort_by", "popularity.desc")
	query.Add("page", "1")
	err := services.CallTMDBApi("/discover/movie", query, &apiResponse)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	c.JSON(200, apiResponse.Results)
}

// @Summary Get movies by genre and release date
// @Description Get movies that match the specified genre and were released in a specific year
// @ID get-movies-by-genre-date
// @Accept json
// @Produce json
// @Param genre query string true "Genre ID"
// @Param year query string false "Release Year"
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /v1/tmdb/movies [get]
func TMDBMoviesByGenreAndDateRouteHandler(c *gin.Context) {
	genre := c.Query("genre")
	year := c.Query("year")
	var movies []types.Movie
	query := url.Values{}
	query.Add("with_genres", genre)
	query.Add("primary_release_year", year)
	err := services.CallTMDBApi("/discover/movie", query, &movies)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	c.JSON(200, movies)
}
