package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	InvalidCookie APIError = 1 << iota
	InvalidState
	ExpiredCookie
	ErrAuthentication
)

type APIError uint16

func (e APIError) Error() string {
	return fmt.Sprintf("Api Error '%s'", strconv.FormatUint(uint64(e), 10))
}

func ConfigCookie(cookie *http.Cookie) *http.Cookie {
	// encode the cookie value using baser64
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	return cookie
}

func SetCookie(c echo.Context, cookie *http.Cookie) (err error) {
	ConfigCookie(cookie)

	// check the length of the cookie
	if len(cookie.String()) > 4096 {
		return fmt.Errorf("%w: coolkie value too long", InvalidCookie)
	}

	c.SetCookie(cookie)

	return
}

func GetCookie(c echo.Context, name string) (*http.Cookie, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return nil, fmt.Errorf("%w: the waiting time for authentication has expired", ExpiredCookie)
		default:
			return nil, fmt.Errorf("%w: an error ocurred while authenticating", ErrAuthentication)
		}
	}

	return cookie, nil
}
