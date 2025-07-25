package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/util"
)

// Havent added cognito field
type AuthValidator struct {
	next    http.Handler
	problem ErrorHandler
	log     *util.CustomLogger
}

func (a *AuthValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.log.Debug().Msg("What are you doing there mate")
	manager, err := key.GetKeyManager()
	if err != nil {
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	authHeaderParts, err := retrieveAuthHeader(r)
	if err != nil {
		a.log.Debug().Msg(err.Error())
		a.problem.HandleError(NewError(w, r, err))
		return
	}

	tokenString := authHeaderParts[1]

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

func retrieveAuthHeader(r *http.Request) ([]string, error) {

	authHeader := r.Header.Get("Authorization")
	splitAuthHeader := strings.Split(authHeader, " ")

	if len(splitAuthHeader) != 2 {
		return nil, fmt.Errorf("No Auth Header")
	}

	return splitAuthHeader, nil
}

func checkRefresh(ctx context.Context, token *key.CustomClaims) (bool, error) {
	if token.Tpe == key.RefreshToken {
		return true, nil
	} else if token.Tpe == key.AccessToken {
		return false, nil
	}
	return false, fmt.Errorf("NO TOKEN TYPE outlined for some flippin reason")
}

func storeClaims(ctx context.Context, token *key.CustomClaims) (context.Context, error) {
	ctx = context.WithValue(ctx, "token", token)
	return ctx, nil
}
