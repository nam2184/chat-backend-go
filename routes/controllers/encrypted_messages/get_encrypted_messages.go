package encrypted_messages

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	queries "github.com/nam2184/generic-queries"
	"github.com/nam2184/mymy/middleware"
	"github.com/nam2184/mymy/models/body"
	qm "github.com/nam2184/mymy/models/db"
	"github.com/nam2184/mymy/models/response"
	"github.com/nam2184/mymy/routes/controllers/options"
	"github.com/nam2184/mymy/routes/controllers/util"
	u "github.com/nam2184/mymy/util"
)

func GetEncryptedMessages(w http.ResponseWriter, r *http.Request, db *sqlx.DB, opts *options.HandlerOptions) {
	ctx := r.Context()

	id, err := getIDFromPath(r)
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
	query_params = util.ValidateQueries[qm.EncryptedMessage](query_params)
	tx := db.MustBegin()

	defer func() {
		if err := tx.Rollback(); err != nil {
			//opts.Problem.HandleError(middleware.NewError(w, r, err))
			return
		}
	}()

	constraint, err := queries.NewConstraint("chat_id = $%", id)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	q, err := queries.SelectOffsetQuery[qm.EncryptedMessage](tx, limit, skip, sort_by, order, query_params, constraint, nil)
	if err != nil {
		opts.Problem.HandleError(middleware.NewError(w, r, err))
		return
	}

	var json_res []byte

	metadata := response.NewMeta()
	var preprocessedArray []body.TempMessage
	for _, message := range q.Rows {
		result := encodeBase64(message)
		preprocessedArray = append(preprocessedArray, result)
	}
	metadata.WithTotal(q.Total)
	resp := response.StandardSuccessResponse(*metadata, preprocessedArray)
	json_res, err = json.Marshal(resp)
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

func getIDFromPath(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	chatIDstr := vars["chatID"]
	chatID, err := u.AtoiToUint(chatIDstr)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func encodeBase64(message qm.EncryptedMessage) body.TempMessage {
	var temp body.TempMessage

	temp.ID = message.ID
	temp.ChatID = message.ChatID
	temp.SenderID = message.SenderID
	temp.SenderName = message.SenderName
	temp.ReceiverID = message.ReceiverID
	temp.Type = message.Type
	temp.IsTyping = message.IsTyping
	temp.Timestamp = message.Timestamp
	temp.Content = message.Content

	if message.Image != nil {
		temp.Image = "data:image/png;base64," + base64.StdEncoding.EncodeToString(message.Image)
	}
	return temp
}
