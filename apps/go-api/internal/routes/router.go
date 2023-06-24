package routes

import (
	"github.com/bnema/flem/go-api/internal/middlewares"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gin-gonic/gin"
)

func SetupRouter(app *types.App) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.SessionMiddleware(app))
	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		user.Use(middlewares.VerifyToken(app))
		{
			movie := user.Group("/movie")
			{
				movie.POST("/preferences", PostUserMoviePreferencesRouteHandler(app)) // posts user movie preferences
				movie.GET("/preferences", GetUserMoviePreferencesRouteHandler(app))   // gets user movie preferences
			}
		}

		whoAmIRoute := v1.Group("/whoami")
		whoAmIRoute.Use(middlewares.VerifyToken(app))
		{
			whoAmIRoute.GET("", WhoAmIRouteHandler(app))
		}
		movieRoute := v1.Group("/movies")
		movieRoute.Use(middlewares.VerifyToken(app))
		{
			movieRoute.GET("", ListMoviesCollectionRouteHandler(app))
		}

		v1.GET("/login", func(c *gin.Context) {
			LoginRoute(app, c)
		})
		v1.GET("/oauth-redirect", func(c *gin.Context) {
			RedirectRoute(app, c)
		})
		tmdb := v1.Group("/tmdb")
		// tmdb.Use(middlewares.IsLoggedIn(app), middlewares.VerifyToken(app))
		{
			tmdb.POST("/movies/post/title", TMDBMovieByTitleRouteHandler)
			tmdb.POST("/movies/post/ids", TMDBMoviesByIDSRouteHandler(app))
			tmdb.GET("/movies", TMDBMoviesByGenreAndDateRouteHandler)
			tmdb.GET("/movies/random10", TMDBRandomMoviesRouteHandler)
		}
		openai := v1.Group("/openai")
		openai.Use(middlewares.VerifyToken(app))
		{
			openai.POST("/movies", SuggestMoviesFromGPT3RouteHandler(app))
			openai.POST("/translate", TranslateMoviesFromGPT3RouteHandler(app))
		}

	}
	return r
}
