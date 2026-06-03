package model

import "errors"

type User struct {
	ID       int64
	ChatID   int64
	Language Language
}

var ErrUserNotFound = errors.New("user not found")
