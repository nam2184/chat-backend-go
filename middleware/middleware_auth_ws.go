package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/util"
)

// Havent added cognito field
type AuthValidatorWS struct {
	next    http.Handler
	problem ErrorHandler
	log     *util.CustomLogger
}

func (a *AuthValidatorWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.log.Debug().Msg("What are you doing there mate")
	manager, err := key.GetKeyManager()
	if err != nil {
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	authParts, err := retrieveAuthToken(r)
	if err != nil {
		a.log.Debug().Msg(err.Error())
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	tokenString := authParts[0]

	a.log.Info().Msgf("Checking token %s", tokenString)
	claims, err := manager.VerifyToken(tokenString)
	if err != nil {
		a.log.Debug().Msg(err.Error())
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	ctx, err := storeClaims(r.Context(), claims)
	if err != nil {
		a.log.Debug().Msg(err.Error())
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	isRefresh, err := checkRefresh(ctx, claims)
	if err != nil {
		a.log.Debug().Msg(err.Error())
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	if isRefresh {
		ctx = context.WithValue(ctx, "refresh", "true")
	}

	a.log.Debug().Msgf("Username :  %s\n", claims.Username)

	r = r.WithContext(ctx)
	a.next.ServeHTTP(w, r)
}

func retrieveAuthToken(r *http.Request) ([]string, error) {
	query := r.URL.Query()
	if len(query) == 0 {
		return nil, fmt.Errorf("No query params found")
	}
	return query["token"], nil
}
