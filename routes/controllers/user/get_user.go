package user

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
	"github.com/nam2184/mymy/routes/controllers/auth/key"
	"github.com/nam2184/mymy/routes/controllers/options"
)

func GetUser(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
    ctx := r.Context()
    ctx = context.WithValue(ctx, "Section", "User get request")
    r = r.WithContext(ctx)
    id, err := getIDFromToken(ctx); if err != nil {
      opts.Problem.HandleError(middleware.NewError(w, r, err))
    }
     
    //Begin inserts with transaction  
    tx := db.MustBegin()
    
    defer func() {
        if err := tx.Rollback(); err != nil {
          //opts.Problem.HandleError(middleware.NewError(w, r, err))
          return
        }
    }()
    
    constraint, err := queries.NewConstraint("id = $%", id); if err != nil {
      opts.Problem.HandleError(middleware.NewError(w, r, err))
      return
    }                           
    
    q, err := queries.SelectQuery[qm.User](tx, constraint, nil)
    if err != nil {
      opts.Problem.HandleError(middleware.NewError(w, r, err))
      return
    }
    //Define response depending on queried user
    var json_res []byte
    if len(q.Rows) == 0 {
      opts.Problem.HandleError(middleware.NewError(w,r, fmt.Errorf("No user found")))
      return
    } else {
	    json_res, err = json.Marshal(q.Rows[0])
      if err != nil {
        opts.Problem.HandleError(middleware.NewError(w, r, err))
        return
      }
    }

	  // Set the Content-Type header to application/json
	  w.Header().Set("Content-Type", "application/json")
	  // Write the JSON response to the HTTP response writer
    w.WriteHeader(http.StatusOK)
    w.Write(json_res)
}


func getIDFromToken(ctx context.Context) (uint, error) {
    token, ok := ctx.Value("token").(*key.CustomClaims); if !ok {
        return 0, fmt.Errorf("No path params found in context")
    }
    id, err := strconv.Atoi(token.Subject); if err != nil {
      return 0, nil
    }
    return uint(id), nil
}


