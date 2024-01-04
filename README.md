# GoError v1.0.0

  
**GoError** is a simple and pluggable Go package that allows developers to create and handle errors in Go more efficiently and gracefully, while also preserving the context of every error. 


## Why use this package? 

 - Plug and Play (compatible with existing go [errors](https://golang.org/pkg/errors/) package)
 - Allows you to specify error types
 - Easy conversion of error into Map or JSON.
 - Add context to error for passing the human-readable error message to client-side applications.
 - Wraps and unwraps error stack. 

## Documentation

Let's look at some code snippets. But before this, make sure to download the package using

    go get -v github.com/dipaksaini98/goerror
   
   and import the package as

    import  (
	    errors "github.com/dipaksaini98/goerror"
    )

### Create error
Syntax:

	New(message string, error_type Type) error


 #### 1. Create an error without an error type


Example:

    err := errors.New("error message", nil)
    fmt.Println(err)
    >> Error message

> Note: when you don't specify any error type *NoType* is set by default.

#### 2. Create an error with an error type

Example:
	

    err :=  errors.New("error message", &errors.BadRequest)
    fmt.Println(err)
    >> Error message

### Add context to an error

Syntax:

	SetContext(error_object error, key string, value interface{}) error

Example:

    err :=  errors.New("programmer friendly error message", nil)
    err := errors.SetContext(err, "user-email", "email is missing")
---
    err :=  errors.New("programmer friendly error message", nil)
    err := errors.SetContext(err, "email", map[string]string{"email_domain": "invalid email domain", "email_local_part": "invalid local part (must start with `a` alphabet)"})

### Get error context

Syntax:

	GetContext(error_object error) map[string]interface{}

Example:

    ctx :=  errors.GetContext(error)

### Add/update error type

Syntax:

	SetType(error_object error, type *Type) error

Example:

    err :=  errors.SetType(err, &errors.NotFound)
  ---
    var customErrorType errors.Type = "CustomErrorType"
    err :=  errors.SetType(err, &customErrorType)

#### List of predefined error types:

	 1. NoType
	 2. BadRequest
	 3. NotFound
	 4. DBError
	 5. Unauthorized
	 6. PermissionDenied
	 7. SomethingWentWrong

### Get error type

Syntax:

	GetType(error_object error) Type

Example:

    errType :=  errors.GetType(err)

### Create a custom error type

Example:

	var CustomErrorType errors.Type = "CustomErrorType"

### Compare two errors (errors Is)

Syntax:

	errors.Is(error_1 error, error_2 error) bool
	
Example:

	err1 := errors.New("error1", nil)
	err2 := errors.New("error2", nil)
	isSame := errors.Is(err1, err2)
	
	fmt.Println(isSame)
	>>false
	
	err3 := err1
	isSame = errors.Is(err3, err1)
	
	fmt.Println(isSame)
	>>true

### Compare error types (errors As)

Syntax:

	errors.As(error_1 error, error_2 error) bool
	
Example:

	err1 := errors.New("error1", nil)
	err2 := errors.New("error2", nil)
	isSameType := errors.As(err1, err2)
	
	fmt.Println(isSameType)
	>>true 
	// Because of NoType (default)
	
	err3 := errors.New("error2", &errors.BadRequest)
	isSameType = errors.As(err3, err1)
	
	fmt.Println(isSameType)
	>>false

### Convert error to Map

Syntax:

	errors.Map(error_1 error) map[string]interface{}
Example:

	err1 := errors.New("error1", nil)
	errMap := errors.Map(err1)
	fmt.Println(errMap)
	>> map[message:error1 type:NoType]
	err1 = errors.SetContext(err1, "email", "email not found")
	errMap = errors.Map(err1)
	fmt.Println(errMap)
	>> map[field:email message:email not found type:NoType]

### Convert error to JSON

Syntax:

	errors.Map(error_1 error) map[string]interface{}
Example:

	err1 := errors.New("error1", nil)
	errJSON := errors.JSON(err1)
	
	fmt.Println(string(errJSON))
	>> {"message":"error1","type":"NoType"}
	
	err1 = errors.SetContext(err1, "email", "email not found")
	errJSON = errors.JSON(err1)
	
	fmt.Println(string(errJSON))
	>> {"field": "email", "message": "email not found", "type": "NoType"}
