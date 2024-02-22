package api

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	status  int
	message string
}

func NewAPIError(status int, message string) *APIError {
	return &APIError{
		status:  status,
		message: message,
	}
}

func (e APIError) Error() string {
	return e.message
}

func (e APIError) Status() int {
	return e.status
}

func ConfigCookie(cookie *http.Cookie) *http.Cookie {
	// encode the cookie value using baser64
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	return cookie
}

func SetCookie(c echo.Context, cookie *http.Cookie) (err *APIError) {
	ConfigCookie(cookie)

	// check the length of the cookie
	if len(cookie.String()) > 4096 {
		err = NewAPIError(http.StatusInternalServerError, "coolkie value too long")
		return
	}

	c.SetCookie(cookie)

	return
}

func GetCookie(c echo.Context, name string) (*http.Cookie, *APIError) {
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
