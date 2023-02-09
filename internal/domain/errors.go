package domain

import "fmt"

var (
	ErrNotFound        = fmt.Errorf("not found")
	ErrRefTypeNotFound = fmt.Errorf("reference type %w", ErrNotFound)
)
