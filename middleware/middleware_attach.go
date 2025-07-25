package middleware

import (
	"context"
	"github.com/nam2184/mymy/util"
	"fmt"
	"net/http"

)

type AttachInfo struct {
  next        http.Handler
  problem     ErrorHandler
  log         *util.CustomLogger 
}

func (a *AttachInfo) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
   
    ctx, err := extractParams(r, ctx); if err != nil {
      a.log.Info().Msg(err.Error())
    }
    
    a.log.Info().Msg("What are you doing attaching this shit bruh")
    r = r.WithContext(ctx)
    a.next.ServeHTTP(w, r)
}


func extractParams(r *http.Request, ctx context.Context) (context.Context, error) {
    query := r.URL.Query()
    if len(query) == 0 {
        return ctx, fmt.Errorf("No query params found")
    }

    params := make(map[string]interface{})
    for key, values := range query {
        if len(values) > 0 {
            params[key] = values[0]
        }
    }

    ctx = context.WithValue(ctx, "query_params", params)
    return ctx, nil
}
