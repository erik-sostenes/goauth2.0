// handlers package is the presentation layer that defines the HTTP handlers
package handlers

import "github.com/labstack/echo/v4"

func GoogleRoutes(
	e *echo.Group,
	googleLoginOAuthHandler GoogleLoginOAuthHandler,
	googleCallbackOAuthHandler GoogleCallbackOAuthHandler,
	m ...echo.MiddlewareFunc,
) {
	e.GET("/login", googleLoginOAuthHandler.Login, m...)
	e.GET("/callback", googleCallbackOAuthHandler.Callback, m...)
}
