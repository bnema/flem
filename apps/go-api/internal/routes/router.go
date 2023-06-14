// routes/router.go
package routes

import (
	"github.com/bnema/flem/go-api/internal/handlers"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *types.App) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
		// Add the LoginRoute and RedirectRoute handlers
		v1.GET("/login", func(c *gin.Context) {
			LoginRoute(app, c)
		})
		v1.GET("/oauth-redirect", func(c *gin.Context) {
			RedirectRoute(app, c)
		})
		v1.POST("/tmdb/movies/post/title", handlers.HandleMoviesByTitle)
		v1.POST("/tmdb/movies/post/ids", handlers.HandleMoviesByIds)
		v1.GET("/tmdb/movies", handlers.HandleMoviesByGenreAndDate)
		v1.GET("/tmdb/movies/random10", handlers.HandleRandomMovies)

	}
	return r
}
