package controllers

import "errors"

var (
	errPayload           = errors.New("invalid data in request body")
	errMissingPathParam  = errors.New("please check for missing path parameter")
	errMissingQueryParam = errors.New("please check for missing query parameter")
	errInvalidPathParam  = errors.New("invalid path parameter")
	errInvalidQueryParam = errors.New("invalid query parameter")
)
