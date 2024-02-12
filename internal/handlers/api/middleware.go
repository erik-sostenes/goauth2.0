package api

import (
	"errors"
	"net/http"

	"github.com/erik-sostenes/auth-api/internal/models"
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

func ErrorHandler(handler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		as := models.UserError(0)

		if !errors.As(err, &as) {
			handler(err, c)
			return
		}

		message := echo.Map{
			"code":    as.Error(),
			"message": err.Error(),
		}

		if errors.As(err, &as) {
			switch as {
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
}