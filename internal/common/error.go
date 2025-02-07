package common

import "net/http"

type ClientError interface {
	HTTPStatus() int
	Code() int
	Error() string
}

type AuthenticationError struct {
	Message string
}

func (e AuthenticationError) HTTPStatus() int {
	return http.StatusUnauthorized
}

func (e AuthenticationError) Code() int {
	return 101
}

func (e AuthenticationError) Error() string {
	return e.Message
}

type AuthorizationError struct {
	Message string
}

func (e AuthorizationError) HTTPStatus() int {
	return http.StatusForbidden
}

func (e AuthorizationError) Code() int {
	return 103
}

func (e AuthorizationError) Error() string {
	return e.Message
}

type InvariantError struct {
	Message string
}

func (e InvariantError) HTTPStatus() int {
	return http.StatusBadRequest
}

func (e InvariantError) Code() int {
	return 100
}

func (e InvariantError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

func (e NotFoundError) Code() int {
	return 104
}

func (e NotFoundError) Error() string {
	return e.Message
}
