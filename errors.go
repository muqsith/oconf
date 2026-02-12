package oconf

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidJSON is returned when the configuration file contains invalid JSON
	ErrInvalidJSON = errors.New("invalid JSON in configuration file")

	// ErrInvalidType is returned when a value has an unexpected type
	ErrInvalidType = errors.New("invalid type in configuration")

	// ErrFileRead is returned when a configuration file cannot be read
	ErrFileRead = errors.New("failed to read configuration file")

	// ErrCircularInclude is returned when circular file inclusion is detected
	ErrCircularInclude = errors.New("circular include detected")
)

// TypeError represents a type assertion error with context
type TypeError struct {
	Key          string
	Expected     string
	Actual       string
	OriginalErr  error
}

func (e *TypeError) Error() string {
	if e.Key != "" {
		return fmt.Sprintf("type error at key '%s': expected %s, got %s", e.Key, e.Expected, e.Actual)
	}
	return fmt.Sprintf("type error: expected %s, got %s", e.Expected, e.Actual)
}

func (e *TypeError) Unwrap() error {
	return e.OriginalErr
}

// FileError represents a file operation error
type FileError struct {
	Path string
	Op   string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file error: %s '%s': %v", e.Op, e.Path, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}
