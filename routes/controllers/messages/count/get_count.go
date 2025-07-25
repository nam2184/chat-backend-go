package count

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/nam2184/mymy/middleware"
	"github.com/nam2184/mymy/routes/controllers/options"
	u "github.com/nam2184/mymy/util"
)

func GetMessagesCount(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
    id, err := getIDFromPath(r); if err != nil {
       opts.Problem.HandleError(middleware.NewError(w, r, err))
       return
    }

    tx := db.MustBegin()
    
    defer func() {
        if err := tx.Rollback(); err != nil {
          //opts.Problem.HandleError(middleware.NewError(w, r, err))
          return
        }
    }()
    
    count, err := getTotalCount(tx, id); if err != nil {
        opts.Problem.HandleError(middleware.NewError(w, r, err))
        return
    }
    
    var json_res []byte
    data := map[string]int{
                "count": count,
            }
    json_res, err = json.Marshal(data)
    if err != nil {
        // Handle the error if JSON marshaling fails
        opts.Problem.HandleError(middleware.NewError(w, r, err))
        return
    } 
    
    // Set the Content-Type header to application/json
	  w.Header().Set("Content-Type", "application/json")
	  // Write the JSON response to the HTTP response writer
    w.WriteHeader(http.StatusOK)
    w.Write(json_res)
}

func getIDFromPath(r *http.Request) (uint, error) {
    vars := mux.Vars(r)
	  chatIDstr := vars["chatID"]
    chatID, err := u.AtoiToUint(chatIDstr);  if err != nil {
        return 0, err
    }
    return chatID, nil
}

func getTotalCount(tx *sqlx.Tx, chatID uint) (int, error) {
    var totalCount int
    countQuery := fmt.Sprintf(
        "SELECT COUNT(*) FROM messages WHERE chat_id = ?",
    )
    err := tx.Get(&totalCount, countQuery, chatID)
    if err != nil {
        return 0, fmt.Errorf("failed to get total count: %w", err)
    }
    
    return totalCount, nil
}


