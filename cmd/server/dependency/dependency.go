// dependency package load the dependencies
package bootstrap

import (
	"os"

	"github.com/erik-sostenes/auth-api/internal/business"
	gh "github.com/erik-sostenes/auth-api/internal/handlers"
	"github.com/erik-sostenes/auth-api/internal/handlers/api"
	"github.com/erik-sostenes/auth-api/internal/models"
	"github.com/erik-sostenes/auth-api/internal/repository"
	"github.com/erik-sostenes/auth-api/internal/repository/persistence"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Injector(e *echo.Echo) (err error) {
	privateKey := os.Getenv("PRIVATE_KEY")
	group := e.Group("/api/v1/goauth/auth")

	e.HTTPErrorHandler = api.ErrorHandler(e.DefaultHTTPErrorHandler)

	oauth := setUpOauth()

	pageDrawer := repository.NewPageDrawer()
	pageProvider := business.NewPageProvider(oauth, pageDrawer)

	set := persistence.NewSet[string, *models.User]()

	codeExchanger := repository.NewCodeExchanger(oauth)
	userInfoAsker := repository.NewUserInfoAsker()
	userSaver := persistence.NewUserSaver(&set)
	generatorToken := business.NewTokenGenerator(privateKey)
	exchanger := business.NewExchanger(codeExchanger, userSaver, userInfoAsker, generatorToken)

	googleHandler := gh.NewGoogleOAuthHandler(pageProvider, exchanger)
	gh.GoogleRoutes(group, googleHandler)

	return
}

func setUpOauth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
