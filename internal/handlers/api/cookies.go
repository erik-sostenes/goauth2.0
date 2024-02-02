package api

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetCookie(c echo.Context, cookie *http.Cookie) error {
	// encode the cookie value using baser64
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// check the length of the cookie
	if len(cookie.String()) > 4096 {
		return errors.New("cookie value too long")
	}

	c.SetCookie(cookie)

	return nil
}

func ReadCookie(c echo.Context, name string) (*http.Cookie, *APIError) {
	cookie, err := c.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, NewAPIError(http.StatusBadRequest, "the waiting time for authentication has expired")
		default:
			return nil, NewAPIError(http.StatusInternalServerError, "an error ocurred while authenticating")
		}
	}

	return cookie, nil
}
