package main

import (
	"fmt"
	"net/url"
	"os"

	docs "github.com/bnema/flem/go-api/docs"
	"github.com/bnema/flem/go-api/internal/routes"
	"github.com/bnema/flem/go-api/pkg/types"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewApp creates a new App struct with all the required fields
func NewApp() *types.App {
	baseUrl, err := url.Parse(os.Getenv("PB_URL"))
	if err != nil {
		panic(err)
	}
	fmt.Println("PB_URL", baseUrl.String())

	app := &types.App{
		SessionStore: sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))),
		PBUrl:        baseUrl.String(),
	}
	fmt.Println("SESSION_SECRET", os.Getenv("SESSION_SECRET"))

	authMethodsUrl, err := baseUrl.Parse("/api/collections/users/auth-methods")
	if err != nil {
		panic(err)
	}
	fmt.Println("authMethodsUrl", authMethodsUrl.String())

	app.PBAuthMethodsURL = authMethodsUrl.String()

	authRefreshUrl, err := baseUrl.Parse("/api/collections/users/auth-refresh")
	if err != nil {
		panic(err)
	}
	fmt.Println("authRefreshUrl", authRefreshUrl.String())
	app.PBAuthRefreshURL = authRefreshUrl.String()

	tradeUrl, err := baseUrl.Parse("/api/collections/users/auth-with-oauth2")
	if err != nil {
		panic(err)
	}
	fmt.Println("tradeUrl", tradeUrl.String())
	app.PBTradeURL = tradeUrl.String()

	oauthRedirectURL, ok := os.LookupEnv("OAUTH_REDIRECT_URL")
	if !ok {
		panic("OAUTH_REDIRECT_URL environment variable is not set")
	}
	app.OAuthRedirectURL = oauthRedirectURL

	PBUserURL, err := baseUrl.Parse("/api/collections/users/records/")
	if err != nil {
		panic(err)
	}
	app.PBUserURL = PBUserURL.String()

	return app
}

// @BasePath /api/v1
func main() {
	app := NewApp()
	r := routes.SetupRouter(app)
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
