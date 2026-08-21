package models

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrInternal                 = errors.New("internal server error")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrGroupNotFound            = errors.New("group not found")
	ErrAppNotFound              = errors.New("app not found")
	ErrAppAlreadyExists         = errors.New("app already exists")
	ErrRedirectURINotFound      = errors.New("redirect uri not found")
	ErrRedirectURIAlreadyExists = errors.New("redirect uri already exists")
	ErrInvalidCredential        = errors.New("invalid credential")
	ErrGroupAlreadyExists       = errors.New("group already exists")
)
