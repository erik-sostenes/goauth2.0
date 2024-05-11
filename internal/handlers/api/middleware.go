package api

import (
	"errors"
	"net/http"

	"github.com/erik-sostenes/auth-api/internal/models"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(handler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		asUserError := models.UserError(0)

		if errors.As(err, &asUserError) {
			message := echo.Map{
				"code":    asUserError.Error(),
				"message": err.Error(),
			}

			if errors.As(err, &asUserError) {
				switch asUserError {
				case models.MissingUserID,
					models.MissingUserName,
					models.MissingUserEmail,
					models.MissingUserPicture,
					models.UserNotFound:
					_ = c.JSON(http.StatusBadRequest, message)
				default:
					_ = c.JSON(http.StatusInternalServerError, message)
				}
			}
		}

		asGoogleAuthError := models.GoogleAuthError(0)

		if errors.As(err, &asGoogleAuthError) {
			message := echo.Map{
				"code":    asGoogleAuthError.Error(),
				"message": err.Error(),
			}

			if errors.As(err, &asGoogleAuthError) {
				switch asGoogleAuthError {
				case models.ErrDuringCodeExchanger:
					_ = c.JSON(http.StatusUnauthorized, message)
				default:
					_ = c.JSON(http.StatusInternalServerError, message)
				}
			}
		}

		asApiError := APIError(0)

		if errors.As(err, &asApiError) {
			message := echo.Map{
				"code":    asApiError.Error(),
				"message": err.Error(),
			}
			switch asApiError {
			case ExpiredCookie:
				_ = c.JSON(http.StatusBadRequest, message)
			case InvalidState:
				_ = c.JSON(http.StatusUnauthorized, message)
			case InvalidCookie, ErrAuthentication:
				_ = c.JSON(http.StatusInternalServerError, message)
			}
		}
	}
}
