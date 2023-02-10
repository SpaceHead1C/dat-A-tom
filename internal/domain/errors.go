package domain

import "fmt"

var (
	ErrNotFound        = fmt.Errorf("not found")
	ErrRecordNotFound  = fmt.Errorf("record %w", ErrNotFound)
	ErrRefTypeNotFound = fmt.Errorf("reference type %w", ErrNotFound)

	ErrExpected = fmt.Errorf("expected")
)
