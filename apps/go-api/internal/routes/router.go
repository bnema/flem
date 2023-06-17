package routes

import (
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
		tmdb := v1.Group("/tmdb")
		// tmdb.Use(middlewares.IsLoggedIn(app), middlewares.VerifyToken(app)) // Using your middlewares here
		{
			tmdb.POST("/movies/post/title", TMDBMovieByTitleRouteHandler)
			tmdb.POST("/movies/post/ids", TMDBMoviesByIDSRouteHandler)
			tmdb.GET("/movies", TMDBMoviesByGenreAndDateRouteHandler)
			tmdb.GET("/movies/random10", TMDBRandomMoviesRouteHandler)
		}
		openai := v1.Group("/openai")
		{
			openai.POST("/movies", SuggestMoviesFromGPT3RouteHandler(app))
			openai.POST("/translate", TranslateMoviesFromGPT3RouteHandler(app))
		}

	}
	return r
}
