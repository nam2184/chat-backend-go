package refresh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/nam2184/mymy/middleware"
	qm "github.com/nam2184/mymy/models/db"
	"github.com/nam2184/mymy/models/response"
	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/util"
)

func GetRefreshToken(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
	ctx := r.Context()
	valid, err := isRefresh(ctx)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	if !valid {
		opts.Problem.HandleError(middleware.NewError(w, r, fmt.Errorf("not a valid refresh token lol")))
		return
	}

	manager, err := key.GetKeyManager()
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	authDetails, ok := ctx.Value("token").(*key.CustomClaims)
	if !ok {
		opts.Problem.HandleError(middleware.NewError(w, r, fmt.Errorf("cant assert to custom claims")))
		return

	}

	sub, err := authDetails.GetSubject()
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, fmt.Errorf("cant assert to custom claims")))
		return
	}

	userid, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, fmt.Errorf("cant assert to custom claims")))
		return
	}

	accessClaims := key.NewNormalClaims(authDetails.Username, userid)
	accessTokenString, err := manager.IssueToken(accessClaims)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	refreshClaims := key.NewRefreshClaims(authDetails.Username, userid)
	refreshTokenString, err := manager.IssueToken(refreshClaims)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	resp := response.NewAuthResponse(accessTokenString, refreshTokenString, accessClaims.ExpiresAt.Time, util.GetZero[qm.User]())
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}

func isRefresh(ctx context.Context) (bool, error) {
	isRefresh, ok := ctx.Value("refresh").(string)
	if !ok {
		return false, fmt.Errorf("no flippin way i cant get refresh from this mf context")
	}

	return isRefresh == "true", nil
}
