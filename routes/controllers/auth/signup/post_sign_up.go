package signup

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/nam2184/mymy/middleware"
	qm "github.com/nam2184/mymy/models/db"
	"github.com/nam2184/mymy/models/request"
	"github.com/nam2184/mymy/models/response"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/util"
	"golang.org/x/crypto/bcrypt"
)

func PostAuthSignUp(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
	ctx := context.WithValue(r.Context(), "Section", "Auth")
	r = r.WithContext(ctx)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	defer r.Body.Close()

	var authDetails request.PostAuthSignUp
	err = json.Unmarshal(data, &authDetails)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	opts.Log.Debug().Msgf("USER IS %s", authDetails.Username)

	tx := db.MustBegin()

	defer func() {
		if err := tx.Rollback(); err != nil {
		}
	}()

	_, err = tx.NamedExec("INSERT INTO auth (username, password) VALUES (:username, :password)", authDetails.ToAuth())
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	stmt, err := tx.PrepareNamed("INSERT INTO users (username, first_name, surname, email, created_at) VALUES (:username, :first_name, :surname, :email, :created_at)")
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	var user qm.User
	err = stmt.Get(&user, authDetails.ToUser())
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	resp := response.PostSignUpResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		Surname:   user.Surname,
		Username:  user.Username,
		Email:     user.Email,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func hashPassword(username, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
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
