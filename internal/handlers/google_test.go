package handlers

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erik-sostenes/auth-api/internal/business"
	"github.com/erik-sostenes/auth-api/internal/repository"
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
			loginHandler := ts.googleLogin

			e.GET(ts.path, loginHandler.Login)

			rq := ts.request
			resp := httptest.NewRecorder()

			e.NewContext(rq, resp)

			if resp.Code != ts.expectedStatusCode {
				t.Log(resp.Body.String())
				t.Errorf("status code was expected %d, but it was obtained %d", ts.expectedStatusCode, resp.Code)
			}
		})
	}
}
