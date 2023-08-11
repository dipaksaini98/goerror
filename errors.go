package goerror

import (
	"errors"
)

type goError struct {
	errorType     Type
	originalError error
	context       Context
	display       bool
	trace         []error
}

type Context struct {
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
func New(msg string, errorType *Type, display bool, context *Context) error {
	if errorType != nil {
		if context != nil {
			return errorType.new(msg, display, context)
		}
		return errorType.new(msg, display, nil)
	} else {
		if context != nil {
			return NoType.new(msg, display, context)
		}
		return NoType.new(msg, display, nil)
	}
}

// Wrap wraps an error
func Wrap(err error, msg string, errorType *Type, display bool, context *Context) error {
	if errorType != nil {
		return errorType.wrap(err, msg, display, context)
	}
	return NoType.wrap(err, msg, display, context)
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
	ctx := Context{key, value}
	if customErr, ok := err.(*goError); ok {
		customErr.context = ctx
		return customErr
	}
	return &goError{errorType: NoType, originalError: err, context: ctx}
}

// GetContext returns the error context
func GetContext(err error) map[string]interface{} {
	emptyCtx := Context{}
	if customErr, ok := err.(*goError); ok && customErr.context != emptyCtx {
		return map[string]interface{}{"field": customErr.context.Key, "message": customErr.context.Value}
	}
	return nil
}
