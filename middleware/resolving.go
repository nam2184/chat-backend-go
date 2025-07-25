package middleware

import (
	"context"
	"net/http"
)

func ValidatingAuthSchema(options MiddlewareOptions) Middleware {
	return func(next http.Handler) http.Handler {
		// Call the next handler with the new context
		return &resolveAuthValidator{
			a: &AuthValidator{
				next:    next,
				problem: options.problemHandler,
				log:     options.log,
			},
			problem: options.problemHandler,
		}
	}
}

type resolveAuthValidator struct {
	a       *AuthValidator
	problem ErrorHandler
}

func (r *resolveAuthValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "Section", "Auth")
	req.WithContext(ctx)

	r.a.ServeHTTP(w, req)
}

func ValidatingAuthWSSchema(options MiddlewareOptions) Middleware {
	return func(next http.Handler) http.Handler {
		// Call the next handler with the new context
		return &resolveAuthWSValidator{
			a: &AuthValidatorWS{
				next:    next,
				problem: options.problemHandler,
				log:     options.log,
			},
			problem: options.problemHandler,
		}
	}
}

type resolveAuthWSValidator struct {
	a       *AuthValidatorWS
	problem ErrorHandler
}

func (r *resolveAuthWSValidator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "Section", "Auth")
	req.WithContext(ctx)

	r.a.ServeHTTP(w, req)
}

func AttachingHeaders(options MiddlewareOptions) Middleware {
	return func(next http.Handler) http.Handler {
		// Call the next handler with the new context
		return &resolveAttachingHeaders{
			a: &AttachHeaders{
				next:    next,
				problem: options.problemHandler,
				log:     options.log,
			},
		}
	}
}

type resolveAttachingHeaders struct {
	a *AttachHeaders
}

func (r *resolveAttachingHeaders) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "Section", "Attaching headers")

	r.a.ServeHTTP(w, req)
}

func AttachingCognitoMetadata(options MiddlewareOptions) Middleware {
	return func(next http.Handler) http.Handler {
		// Call the next handler with the new context
		return &resolveCognitoMetadata{
			a: &AttachInfo{
				next:    next,
				problem: options.problemHandler,
				log:     options.log,
			},
			problem: options.problemHandler,
		}
	}
}

type resolveCognitoMetadata struct {
	a       *AttachInfo
	problem ErrorHandler
}

func (r *resolveCognitoMetadata) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "Section", "Get cognito metadata")
	ctx = context.WithValue(ctx, "Logger", r.a.log)
	req = req.WithContext(ctx)

	r.a.ServeHTTP(w, req)
}
