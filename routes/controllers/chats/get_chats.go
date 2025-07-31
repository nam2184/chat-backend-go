package chats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	queries "github.com/nam2184/generic-queries"
	"github.com/nam2184/mymy/middleware"
	qm "github.com/nam2184/mymy/models/db"
	"github.com/nam2184/mymy/models/response"
	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/routes/controllers/util"
)

func GetChats(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
	ctx := r.Context()

	id, err := getIDFromToken(ctx)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}
	query_params, ok := ctx.Value("query_params").(map[string]interface{})
	if !ok {
		opts.Log.Debug().Msg("No query params")
		query_params = make(map[string]interface{}, 0)
	}

	skip, limit, err := util.GetSkipLimit(query_params)
	if err != nil {
		opts.Log.Debug().Msg(err.Error())
	}

	opts.Log.Debug().Msgf("Limit %d", limit)
	opts.Log.Debug().Msgf("Skip %d", skip)

	sort_by, order := util.GetSortBy(query_params)

	tx := db.MustBegin()

	defer func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}()

	constraint, err := queries.NewConstraint("user1_id = $% OR user2_id = $%", id, id)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	q, err := queries.SelectOffsetQuery[qm.Chat](tx, limit, skip, sort_by, order, query_params, constraint, nil)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	users, err := getChatUsers(tx, id)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	chatResponse := response.GetChats{
		Chats: q.Rows,
		Users: users,
	}

	metadata := response.NewMeta()
	metadata.WithTotal(q.Total)

	resp := response.StandardSuccessResponse(*metadata, chatResponse)
	json_res, err := json.Marshal(resp)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response to the HTTP response writer
	w.WriteHeader(http.StatusOK)
	w.Write(json_res)
}

func getIDFromToken(ctx context.Context) (int64, error) {
	token, ok := ctx.Value("token").(*key.CustomClaims)
	if !ok {
		return 0, fmt.Errorf("No path params found in context")
	}
	id, err := strconv.ParseInt(token.Subject, 10, 64)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func getChatUsers(db *sqlx.Tx, userID int64) ([]qm.User, error) {
	query := `
		SELECT u.*
		FROM users u
		JOIN (
			SELECT 
				CASE 
					WHEN user1_id = $1 THEN user2_id
					ELSE user1_id 
				END AS other_user_id
			FROM chats
			WHERE $1 IN (user1_id, user2_id)
		) c ON u.id = c.other_user_id
	`
	users := make([]qm.User, 0)
	err := db.Select(&users, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}
