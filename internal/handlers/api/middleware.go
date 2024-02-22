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
				case models.DuplicateUser:
					_ = c.JSON(http.StatusBadRequest, message)
				case models.MissingUserID,
					models.MissingUserName,
					models.MissingUserEmail,
					models.MissingUserPicture:
					_ = c.JSON(http.StatusBadRequest, message)
				default:
					handler(err, c)
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
					handler(err, c)
				}
			}
		}

		handler(err, c)
	}
}
