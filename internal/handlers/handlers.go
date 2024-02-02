// handlers package is the presentation layer that defines the HTTP handlers
package handlers

import "github.com/labstack/echo/v4"

func GoogleRoutes(e *echo.Group, handlers GoogleOAuthHandler, m ...echo.MiddlewareFunc) {
	e.GET("/login", handlers.Login, m...)
	e.GET("/callback", handlers.Callback, m...)
}
