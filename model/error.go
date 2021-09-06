package model

type ErrNotFound struct {
	error
}

type error interface {
	Error() string
}
