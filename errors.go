package goinsta

import "errors"

var (
	// ErrNotFound is returned if the request responds with a 404 status code
	// i.e a non existent user
	ErrNotFound = errors.New("The specified data wasn't found.")

	// ErrLoggedOut is returned if the request responds with a 400 status code
	ErrLoggedOut = errors.New("The account is logged out")

	// ErrNoID
	ErrNoID = errors.New("User id have not being specified. See the documentation at godoc.org")
)
