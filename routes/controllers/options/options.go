package options

import (
	"github.com/nam2184/mymy/middleware"
	"github.com/nam2184/mymy/util"
)

type HandlerOption func(*HandlerOptions)

type HandlerOptions struct {
    Problem        middleware.ErrorHandler
    Log            *util.CustomLogger
}


// WithErrorHandler returns a Handler option that sets Problem handler.
func WithErrorHandler(e middleware.ErrorHandler) HandlerOption {
		return func(c *HandlerOptions) {
      c.Problem = e
    }
}

func WithLogger(logger *util.CustomLogger) HandlerOption {
		return func(c *HandlerOptions) {
      c.Log = logger
    }
}

func NewHandlerConfig(opts ...HandlerOption) *HandlerOptions {
    config := &HandlerOptions{}

    for _, opt := range opts {
        opt(config)
    }

    return config
}
