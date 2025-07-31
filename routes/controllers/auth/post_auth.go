package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nam2184/mymy/middleware"
	qm "github.com/nam2184/mymy/models/db"
	"github.com/nam2184/mymy/models/response"
	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/util"
	"golang.org/x/crypto/bcrypt"
)

func PostAuthenticated(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
	ctx := context.WithValue(r.Context(), "Section", "Auth")
	r = r.WithContext(ctx)

	manager, err := key.GetKeyManager()
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	defer r.Body.Close()

	var authDetails qm.Auth
	err = json.Unmarshal(data, &authDetails)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	opts.Log.Debug().Msgf("USER IS %s", authDetails.Username)

	valid, err := validateCredentials(db, authDetails.Username, authDetails.Password)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	if !valid {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	user, err := retrieveUser(db, authDetails.Username)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	accessCfg := key.TokenConfig{
		Username: authDetails.Username,
		ID:       user.ID,
		Type:     key.AccessToken,
		Expiry:   24 * time.Hour,
		Issuer:   "auth-service",
	}

	accessClaims := key.NewClaims(accessCfg)
	accessTokenString, err := manager.IssueToken(accessClaims)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	refreshCfg := key.TokenConfig{
		Username: authDetails.Username,
		ID:       user.ID,
		Type:     key.RefreshToken,
		Expiry:   72 * time.Hour,
		Issuer:   "auth-service",
	}

	refreshClaims := key.NewClaims(refreshCfg)
	refreshTokenString, err := manager.IssueToken(refreshClaims)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	resp := response.PostAuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		Expiry:       accessClaims.ExpiresAt.Time,
		User:         user,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func validateCredentials(db *sqlx.DB, username, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password FROM auth WHERE username = $1`
	err := db.Get(&hashedPassword, query, username)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}

func retrieveUser(db *sqlx.DB, username string) (qm.User, error) {
	var user qm.User
	query := `SELECT * FROM users WHERE username = $1`
	err := db.Get(&user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return util.GetZero[qm.User](), fmt.Errorf("no user found")
		}
		return util.GetZero[qm.User](), err
	}
	return user, nil
}
