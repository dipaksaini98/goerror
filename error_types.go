package goerror

import "errors"

// Type is the type of an error
type Type string

var (
	// NoType default error type when no error type is specified
	NoType Type = "NoType"
	// BadRequest indicates that the server cannot or will not process the request due to something that is perceived to be a client error
	BadRequest Type = "BadRequest"
	// NotFound indicates that the server can't find the requested resource.
	NotFound Type = "NotFound"
	// DBError indicates that the server can't find the requested resource due to some error occurred while querying the database.
	DBError Type = "DBError"
	//Unauthorized indicates that the request lacks valid authentication credentials for the target resource.
	Unauthorized Type = "Unauthorized"
	// PermissionDenied indicates that client does not have access rights to the resource
	PermissionDenied Type = "PermissionDenied"
	// SomethingWentWrong indicates that server has encountered a situation it doesn't know how to handle
	SomethingWentWrong Type = "SomethingWentWrong"
	// InternalServerError indicates that server has encountered a situation it doesn't know how to handle
	InternalServerError Type = "InternalServerError"
)

// new creates a new custom error object
func (errType Type) new(title string, msg string, display bool) error {
	err := &GoError{OriginalError: errors.New(msg), ErrorType: errType, Display: display, Message: msg, Title: title}
	err.Trace = append(err.Trace, err)
	return err
}

// wrap wraps context with an error object
func (errType Type) wrap(err error, systemErr error, title string, msg string, display bool) error {
	newErr := &GoError{ErrorType: errType, OriginalError: systemErr, Display: display, Message: msg, Title: title}
	err.(*GoError).Trace = append(err.(*GoError).Trace, newErr)
	return err
}
