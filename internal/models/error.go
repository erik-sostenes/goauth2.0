package models

import (
	"strconv"
)

const (
	MissingUserID UserError = 1 << iota
	MissingUserName
	MissingUserEmail
	MissingUserPicture
	DuplicateUser
)

type UserError uint16

func (u UserError) Error() string {
	return "user = " + strconv.FormatUint(uint64(u), 10)
}

const (
	ErrDuringCodeExchanger GoogleAuthError = 1 << iota
)

type GoogleAuthError uint16

func (g GoogleAuthError) Error() string {
	return "google = " + strconv.FormatUint(uint64(g), 10)
}
