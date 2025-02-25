package service

import "errors"

var (
	ErrNotEnoughStock       = errors.New("not enough stock for")
	ErrAlreadyExists        = errors.New("already exists")
	ErrInvalidPasswordEmail = errors.New("invalid password or email")
)
