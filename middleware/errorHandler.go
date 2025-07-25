package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
  "context"
)

// NewError returns a new problem occurred while processing the request.
func NewError(w http.ResponseWriter, req *http.Request, err error) Error {	 
  return Error{
		w:   w,
		req: req,
    err: NewErrorType(req.Context(), err),
	}
}

func NewErrorSocket(conn *websocket.Conn, err error) Error {	 
  return Error{
    conn: conn,
    err: NewErrorType(nil, err),
	}
}



func NewMultiError(message string, errs ...error) MultiError {
  return MultiError{
    message : message,
    err: errs,
  }
}
// Error describes a problem occurred while processing the request (or the response).
// In most cases, the problem represents a validation error.

type Error struct {
	w   http.ResponseWriter
	req *http.Request
	conn *websocket.Conn
	err ErrorType
}

type MultiError struct {
  message string
  err []error
}

// Cause returns the underlying error that represents the problem.
func (p Error) Cause() string {
	return p.err.Message
}

func (p Error) FaultPart() string {
	return p.err.Section
}


// ResponseWriter retruns the ResponseWriter relative to the request.
func (p Error) ResponseWriter() http.ResponseWriter {
  return p.w
}

// Request returns the request on which the problem occured.
func (p Error) Request() *http.Request {
	return p.req
}

// ErrorHandlerFunc is a function that handles problems occurred in a middleware
// while processing a request or a response.
//
// This function implements ErrorHandler.
type ErrorHandlerFunc func(Error)

// HandleError handles the problem.
func (f ErrorHandlerFunc) HandleError(problem Error) {
	f(problem)
}

type ErrorType struct {
	Section string `json:"section"`
	Message string `json:"message"`
	err     error  
}

func NewErrorType(ctx context.Context, err error) ErrorType {
  var section string
  var ok bool
  if ctx != nil {
    section, ok = ctx.Value("Section").(string); if !ok {
      return ErrorType{ Message: err.Error(), err: err} 
    }
  } else {
    section = "Unknown, prob websockky "
  }

  return ErrorType{Section : section, Message: err.Error(), err : err,}
}

func AddError(err1, err2 error) error {
    if err1 == nil && err2 == nil {
        return nil
    }
    if err1 == nil {
        return err2
    }
    if err2 == nil {
        return err1
    }

    return fmt.Errorf("%s; %s", err1.Error(), err2.Error())
}


type ErrorResponse struct {
	  Status     string          `json:"status"`
	  Type       int             `json:"type"`
	  Message    string          `json:"message"`
	  Code       int             `json:"code"`
	  Error      ErrorType       `json:"error"`
}

func StandardErrorResponse(error ErrorType, code int) *ErrorResponse {
  return &ErrorResponse{ Status: "failure", Type: code, Error: error}
}

// ErrorHandler can handle problems occurred in a middleware while processing
// a request or a response.
//
// Validation error depends on the middleware type, e.g. for query validator
// middleware the error will describe query validation failure. Usually, the
// handler should not wrap the error with a message like "query validation failure",
// because the message will be already present in such error.
type ErrorHandler interface {
	HandleError(problem Error)
}

// newErrorHandlerErrorResponder is a very simple ErrorHandler that
// writes problem error message to the response.

func NewErrorHandlerErrorResponder() ErrorHandlerFunc {
	return func(p Error) {
    p.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Println(p.err.Section)
    var code int
    if p.err.Section == "Auth"  {
      code = http.StatusUnauthorized 
    } else {
      code = http.StatusBadRequest 
    }

    p.ResponseWriter().WriteHeader(code)
    response := StandardErrorResponse(p.err, code)
    err_json, _ := json.Marshal(response)
    p.ResponseWriter().Write([]byte(err_json)) 
	}
}

// newErrorHandlerWarnLogger is a very simple ErrorHandler that writes
// problem error to the standard logger with a warning prefix.
func NewErrorHandlerWarnLogger(kind string) ErrorHandlerFunc {
	return func(p Error) {
		log.Printf("[WARN]  %s problem on \"%s %s\": %v", kind, p.Request().Method, p.Request().URL.String(), p.Cause())
	}
}


