package goerror

import (
	"errors"
)

type goError struct {
	errorType     Type
	originalError error
	title         string
	message       string
	context       context
	display       bool
	trace         []error
}

type context struct {
	Key   interface{}
	Value interface{}
}

// Error implements the Error interface required by Go
func (e *goError) Error() string {
	return e.originalError.Error()
}

// Unwrap implements the Unwrap interface required by Go
func (e *goError) Unwrap() error {
	return errors.Unwrap(e.originalError)
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
	if goErr, ok := err.(*goError); ok {
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
	if goErr, ok := err.(*goError); ok {
		if targetGoErr, ok1 := target.(*goError); ok1 {
			return goErr.errorType == targetGoErr.errorType
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
	if customErr, ok := err.(*goError); ok {
		customErr.errorType = t
		return customErr
	}
	return &goError{errorType: t, originalError: err}
}

// GetType returns the error type
func GetType(err error) Type {
	if goErr, ok := err.(*goError); ok {
		return goErr.errorType
	}
	return NoType
}

// SetContext adds context to the error
func SetContext(err error, key, value interface{}) error {
	ctx := context{key, value}
	if customErr, ok := err.(*goError); ok {
		customErr.context = ctx
		return customErr
	}
	return &goError{errorType: NoType, originalError: err, context: ctx}
}

// GetContext returns the error context
func GetContext(err error) map[string]interface{} {
	emptyCtx := context{}
	if customErr, ok := err.(*goError); ok && customErr.context != emptyCtx {
		return map[string]interface{}{"field": customErr.context.Key, "message": customErr.context.Value}
	}
	return nil
}
