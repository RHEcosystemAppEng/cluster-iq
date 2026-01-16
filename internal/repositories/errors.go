package repositories

import "errors"

// ErrNotFound is returned when a resource is not found in the database.
var ErrNotFound = errors.New("requested resource not found")
