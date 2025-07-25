package middleware

import (
	"net/http"
	"regexp"

	"github.com/getkin/kin-openapi/routers"
	"github.com/nam2184/mymy/util"
)

// Middleware describes a middleware that can be applied to a http.handler.
type Middleware func(next http.Handler) http.Handler

// MiddlewareOptions represent options for middleware.
type MiddlewareOptions struct {
    router          routers.Router
    jsonSelectors   []*regexp.Regexp
    problemHandler  ErrorHandler
    continueOnError bool
    log             *util.CustomLogger
}

func CreateMiddlewareOptions() *MiddlewareOptions {
    return &MiddlewareOptions{}
}
// WithJSONSelectors returns a middleware option that sets JSON Content-Type selectors.
func (m *MiddlewareOptions) WithJSONSelectors(selectors ...*regexp.Regexp) {
		m.jsonSelectors = append(m.jsonSelectors, selectors...)
}

// WithErrorHandler returns a middleware option that sets problem handler.
func (m *MiddlewareOptions) WithErrorHandler(h ErrorHandler)  {
		m.problemHandler = h
}

// WithErrorHandlerFunc returns a middleware option that sets problem handler.
func (m *MiddlewareOptions) WithErrorHandlerFunc(f ErrorHandlerFunc) {  
		m.problemHandler = f
}

// WithContinueOnError returns a middleware option that defines if middleware
// should continue when error occurs.
func (m *MiddlewareOptions) WithContinueOnError(contin bool)  {
		m.continueOnError = contin
}

func (m *MiddlewareOptions) WithLogger(logger *util.CustomLogger)  {
    m.log = logger
}


