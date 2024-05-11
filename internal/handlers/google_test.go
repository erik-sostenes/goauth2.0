package handlers

import (
	"context"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/erik-sostenes/auth-api/internal/business"
	"github.com/erik-sostenes/auth-api/internal/handlers/api"
	"github.com/erik-sostenes/auth-api/internal/models"
	"github.com/erik-sostenes/auth-api/internal/repository"
	"github.com/erik-sostenes/auth-api/internal/repository/memory"
	"github.com/erik-sostenes/auth-api/pkg"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	clientId     = flag.String("client_id", "", "google client id")
	clientSecret = flag.String("client_secret", "", "google client secret")
	redirectUrl  = flag.String("redirect_url", "", "google redirect url")
)

func TestHttpHandler_GoogleLogin(t *testing.T) {
	tsc := map[string]struct {
		googleLogin        GoogleLoginOAuthHandler
		path               string
		request            *http.Request
		expectedStatusCode int
	}{
		"given the correct configuration, a template will be rendered to authenticate with google": {
			googleLogin: func() GoogleLoginOAuthHandler {
				oauth := &oauth2.Config{
					ClientID:     *clientId,
					ClientSecret: *clientSecret,
					RedirectURL:  *redirectUrl,
					Scopes: []string{
						"https://www.googleapis.com/auth/userinfo.email",
						"https://www.googleapis.com/auth/userinfo.profile",
					},
					Endpoint: google.Endpoint,
				}

				pageDrawer := repository.NewPageDrawer()
				pageProvider := business.NewPageProvider(oauth, pageDrawer)
				return NewGoogleLoginOAuthHandler(pageProvider)
			}(),
			path:               "/api/v1/goauth/auth",
			request:            httptest.NewRequest(http.MethodGet, "/api/v1/goauth/auth", http.NoBody),
			expectedStatusCode: http.StatusOK,
		},
	}

	for name, ts := range tsc {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			e.HTTPErrorHandler = api.ErrorHandler(e.DefaultHTTPErrorHandler)

			loginHandler := ts.googleLogin

			e.GET(ts.path, loginHandler.Login)

			rq := ts.request
			resp := httptest.NewRecorder()

			e.ServeHTTP(resp, rq)

			if resp.Code != ts.expectedStatusCode {
				t.Log(resp.Body.String())
				t.Errorf("status code was expected %d, but it was obtained %d", ts.expectedStatusCode, resp.Code)
			}
		})
	}
}

func TestHttpHandler_GoogleCallback(t *testing.T) {
	oauthConfig := &oauth2.Config{
		ClientID:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	tsc := map[string]struct {
		googleCallback     GoogleCallbackOAuthHandler
		request            *http.Request
		cookie             *http.Cookie
		expectedStatusCode int
	}{
		"given the status cookie was not set, a BadRequest status code was expected": {
			googleCallback: func() GoogleCallbackOAuthHandler {
				codeExchanger := repository.NewCodeExchanger(oauthConfig)

				set := pkg.NewSet[string, *models.User]()
				userMemory := memory.NewUserMemory(&set)

				userInfoAsker := repository.NewUserInfoAsker()

				privateKey := os.Getenv("PRIVATE_KEY")
				tokenGenerator := business.NewTokenGenerator(privateKey)

				exchanger := business.NewExchanger(codeExchanger, userMemory, userMemory, userInfoAsker, tokenGenerator)
				return NewGoogleCallbackOAuthHandler(exchanger)
			}(),
			request:            httptest.NewRequest(http.MethodGet, "/api/v1/goauth/auth/callback", http.NoBody),
			expectedStatusCode: http.StatusBadRequest,
		},
		"given the callback url is missing the query parameter 'state', a StatusUnauthorized status code was expected": {
			googleCallback: func() GoogleCallbackOAuthHandler {
				codeExchanger := repository.NewCodeExchanger(oauthConfig)

				set := pkg.NewSet[string, *models.User]()
				userMemory := memory.NewUserMemory(&set)

				userInfoAsker := repository.NewUserInfoAsker()

				privateKey := os.Getenv("PRIVATE_KEY")
				tokenGenerator := business.NewTokenGenerator(privateKey)

				exchanger := business.NewExchanger(codeExchanger, userMemory, userMemory, userInfoAsker, tokenGenerator)
				return NewGoogleCallbackOAuthHandler(exchanger)
			}(),
			request: httptest.NewRequest(http.MethodGet, "/api/v1/goauth/auth/callback", http.NoBody),
			cookie: &http.Cookie{
				Name:     cookieName,
				Value:    "012-345-6789",
				SameSite: http.SameSiteLaxMode,
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		"given the callback url is missing the query parameter 'code', a StatusUnauthorized status code was expected": {
			googleCallback: func() GoogleCallbackOAuthHandler {
				codeExchanger := repository.NewCodeExchanger(oauthConfig)

				set := pkg.NewSet[string, *models.User]()
				userMemory := memory.NewUserMemory(&set)

				userInfoAsker := repository.NewUserInfoAsker()

				privateKey := os.Getenv("PRIVATE_KEY")
				tokenGenerator := business.NewTokenGenerator(privateKey)

				exchanger := business.NewExchanger(codeExchanger, userMemory, userMemory, userInfoAsker, tokenGenerator)
				return NewGoogleCallbackOAuthHandler(exchanger)
			}(),
			request: httptest.NewRequest(http.MethodGet, "/api/v1/goauth/auth/callback?state=012-345-6789", http.NoBody),
			cookie: &http.Cookie{
				Name:     cookieName,
				Value:    "012-345-6789",
				SameSite: http.SameSiteLaxMode,
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		"given a successful authentication, a StatusSeeOther status code was expected": {
			googleCallback: func() GoogleCallbackOAuthHandler {
				codeExchanger := NewCodeExchangerMock(oauthConfig)

				set := pkg.NewSet[string, *models.User]()
				userMemory := memory.NewUserMemory(&set)

				userInfoAsker := NewUserInfoAskerMock()

				privateKey := os.Getenv("PRIVATE_KEY")
				tokenGenerator := business.NewTokenGenerator(privateKey)

				exchanger := business.NewExchanger(codeExchanger, userMemory, userMemory, userInfoAsker, tokenGenerator)
				return NewGoogleCallbackOAuthHandler(exchanger)
			}(),
			request: httptest.NewRequest(http.MethodGet, "/api/v1/goauth/auth/callback?state=012-345-6789&code=4/0AeaYSHA496Vz", http.NoBody),
			cookie: &http.Cookie{
				Name:     cookieName,
				Value:    "012-345-6789",
				SameSite: http.SameSiteLaxMode,
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusSeeOther,
		},
	}

	for name, ts := range tsc {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			e.HTTPErrorHandler = api.ErrorHandler(e.DefaultHTTPErrorHandler)

			e.GET("/api/v1/goauth/auth/callback", ts.googleCallback.Callback)

			rq := ts.request
			resp := httptest.NewRecorder()

			if ts.cookie != nil {
				api.ConfigCookie(ts.cookie)
				rq.AddCookie(ts.cookie)
			}

			e.ServeHTTP(resp, rq)

			if resp.Code != ts.expectedStatusCode {
				t.Log(resp.Body.String())
				t.Errorf("status code was expected %d, but it was obtained %d", ts.expectedStatusCode, resp.Code)
			}
		})
	}
}

type codeExchangerMock struct {
	config *oauth2.Config
}

func NewCodeExchangerMock(config *oauth2.Config) *codeExchangerMock {
	return &codeExchangerMock{
		config: config,
	}
}

func (c *codeExchangerMock) ExchangeCode(ctx context.Context, code string) (models.Token, error) {
	return &oauth2.Token{}, nil
}

type userInfoAskerMock struct{}

func NewUserInfoAskerMock() *userInfoAskerMock {
	return &userInfoAskerMock{}
}

func (userInfoAskerMock) AskUserInfo(ctx context.Context, token models.Token) (user *models.User, err error) {
	return models.NewUser(
		"1",
		"Erik Sostenes Simon",
		"eriksostenessimon@gmail.com",
		"https://eriksostenessimon.com",
		true,
	)
}
