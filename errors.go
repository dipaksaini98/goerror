package goerror

import (
	"errors"
)

type GoError struct {
	ErrorType     Type
	OriginalError error
	Title         string
	Message       string
	Context       Context
	Display       bool
	Trace         []error
}

type Context struct {
	Key   interface{}
	Value interface{}
}

// Error implements the Error interface required by Go
func (e *GoError) Error() string {
	return e.OriginalError.Error()
}

// Unwrap implements the Unwrap interface required by Go
func (e *GoError) Unwrap() error {
	return errors.Unwrap(e.OriginalError)
}

// New returns new error
func New(title string, msg string, errorType *Type, display bool) error {
	if errorType != nil {
		return errorType.new(title, msg, display)
	} else {
		return NoType.new(title, msg, display)
	}
}

// Wrap wraps an error
func Wrap(err error, systemErr error, title string, msg string, errorType *Type, display bool) error {
	if errorType != nil {
		return errorType.wrap(err, systemErr, title, msg, display)
	}
	return NoType.wrap(err, systemErr, title, msg, display)
}

// Unwrap unwraps an error
func Unwrap(err error) error {
	if goErr, ok := err.(*GoError); ok {
		return goErr.Unwrap()
	}
	return errors.Unwrap(err)
}

// Is implements the Is interface defined in go specification
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As implements the As interface defined in go specification
func As(err error, target error) bool {
	if goErr, ok := err.(*GoError); ok {
		if targetGoErr, ok1 := target.(*GoError); ok1 {
			return goErr.ErrorType == targetGoErr.ErrorType
		}
		return false
	}
	return false
}

// Error returns the error message
func Error(err error) string {
	return err.Error()
}

// SetType add/change the error type
func SetType(err error, t Type) error {
	if customErr, ok := err.(*GoError); ok {
		customErr.ErrorType = t
		return customErr
	}
	return &GoError{ErrorType: t, OriginalError: err}
}

// GetType returns the error type
func GetType(err error) Type {
	if goErr, ok := err.(*GoError); ok {
		return goErr.ErrorType
	}
	return NoType
}

// GetTitle returns the error title
func GetTitle(err error) string {
	if goErr, ok := err.(*GoError); ok {
		return goErr.Title
	}
	return ""
}

// GetDisplay returns whether the error should be displayed
func GetDisplay(err error) bool {
	if goErr, ok := err.(*GoError); ok {
		return goErr.Display
	}
	return false
}

// GetTrace returns the error trace
func GetTrace(err error) []error {
	if goErr, ok := err.(*GoError); ok {
		return goErr.Trace
	}
	return nil
}

// GetOriginalError returns the original error
func GetOriginalError(err error) error {
	if goErr, ok := err.(*GoError); ok {
		return goErr.OriginalError
	}
	return err
}

// SetContext adds context to the error
func SetContext(err error, key, value interface{}) error {
	ctx := Context{key, value}
	if customErr, ok := err.(*GoError); ok {
		customErr.Context = ctx
		return customErr
	}
	return &GoError{ErrorType: NoType, OriginalError: err, Context: ctx}
}

// GetContext returns the error context
func GetContext(err error) map[string]interface{} {
	emptyCtx := Context{}
	if customErr, ok := err.(*GoError); ok && customErr.Context != emptyCtx {
		return map[string]interface{}{"field": customErr.Context.Key, "message": customErr.Context.Value}
	}
	return nil
}
