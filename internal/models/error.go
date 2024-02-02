package models

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	MissingUserID UserError = 1 << iota
	MissingUserName
	MissingUserEmail
	MissingUserPicture
)

type UserError uint16

func (u UserError) Error() string {
	return u.String()
}

func (u UserError) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

func (UserError) Join(err ...error) error {
	return errors.Join(err...)
}

func (u UserError) Wrap(userErr UserError, msg string) error {
	return fmt.Errorf("%w: %s", userErr, msg)
}
