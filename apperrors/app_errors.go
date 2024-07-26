package apperrors

import "net/http"

type AppError interface {
	Error() string
	Code() int
}

type NotFound struct {
	Message string
}

func (e NotFound) Error() string {
	return e.Message
}

func (e NotFound) Code() int {
	return http.StatusNotFound
}

type Validation struct {
	Message string
}

func (e Validation) Error() string {
	return e.Message
}

func (e Validation) Code() int {
	return http.StatusBadRequest
}

type Unauthorized struct {
	Message string
}

func (e Unauthorized) Error() string {
	return e.Message
}

func (e Unauthorized) Code() int {
	return http.StatusUnauthorized
}

type Internal struct {
	Message string
}

func (e Internal) Error() string {
	return e.Message
}

func (e Internal) Code() int {
	return http.StatusInternalServerError
}
