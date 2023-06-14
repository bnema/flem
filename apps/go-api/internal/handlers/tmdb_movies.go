package handlers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/bnema/flem/go-api/internal/services"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
)

// @Summary Search movies by title
// @Description Get movies that match given titles
// @ID get-movies-by-title
// @Accept json
// @Produce json
// @Param titles body []string true "List of Titles"
// @Success 200 {array} types.Movie
// @Failure 500 {object} types.Error
// @Router /v1/tmdb/movies/post/title [post]
func HandleMoviesByTitle(c *gin.Context) {
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
// @Router /v1/tmdb/movies/post/ids [post]
func HandleMoviesByIds(c *gin.Context) {
	idsStr := c.PostFormArray("ids")
	ids := make([]int, len(idsStr))
	for i, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Something went wrong",
			})
			return
		}
		ids[i] = id
		var movie types.Movie
		err = services.CallTMDBApi(fmt.Sprintf("/movie/%d", id), url.Values{}, &movie)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Something went wrong",
			})
			return
		}
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
// @Router /v1/tmdb/random10 [get]
func HandleRandomMovies(c *gin.Context) {
	var movies []types.Movie
	query := url.Values{}
	query.Add("sort_by", "popularity.desc")
	query.Add("page", "1")
	err := services.CallTMDBApi("/discover/movie", query, &movies)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	c.JSON(200, movies)
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
func HandleMoviesByGenreAndDate(c *gin.Context) {
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
