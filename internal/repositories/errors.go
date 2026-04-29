package repositories

import "errors"

// ErrNotFound is returned when a resource is not found in the database.
var ErrNotFound = errors.New("requested resource not found")

// ErrNoClustersInAccount is returned when an account exists but has no associated clusters.
var ErrNoClustersInAccount = errors.New("no clusters found for this account")
