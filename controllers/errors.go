package controllers

import "errors"

var (
	errPayload          = errors.New("invalid data in request body")
	ErrMissingPathParam = errors.New("please check for missing path parameter")
	ErrInvalidPathParam = errors.New("invalid path parameter")
)
