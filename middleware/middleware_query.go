package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/nam2184/mymy/util"
)


type SchemaValidator struct {
  next                http.Handler
  router              routers.Router
  doc                 *openapi3.T
  problem             ErrorHandler
  log                 *util.CustomLogger 
  continueOnProblem   bool
}

func  (s *SchemaValidator) ServeHTTP(w http.ResponseWriter, r *http.Request, loader *openapi3.Loader)  {
      route, pathParams, err := s.router.FindRoute(r)
      s.log.Info().Msg("Starting schema validation \n")
      if err != nil {
        s.problem.HandleError(NewError(w,r, err))
        return
      }
      
      requestValidationInput := &openapi3filter.RequestValidationInput  {
					Request:    r,
					PathParams: pathParams,
          QueryParams: r.URL.Query(),
					Route:      route,
			}
      
      s.log.Debug().Msgf("Validating request :  %+v\n", requestValidationInput)
      
      if err := openapi3filter.ValidateRequest(loader.Context, requestValidationInput); err != nil { 
          s.problem.HandleError(NewError(w, r, err))
          return
      }
      
      ctx, err := getPathParams(r.Context(), pathParams); if err != nil {
          s.log.Info().Msg(err.Error())
      }

      r = r.WithContext(ctx)
      s.log.Info().Msg("Finished validating schema \n")

      s.log.Debug().Msgf("Next handler :  %+v\n", s.next)
      s.next.ServeHTTP(w, r)
}

func getPathParams(ctx context.Context, paths map[string]string) (context.Context, error) {
    if len(paths) == 0 {
      return ctx, fmt.Errorf("No path params")
    } 
    ctx = context.WithValue(ctx, "path_params", paths)
    return ctx, nil
}

func removeV1(endpoint string) string {
    parts := strings.Split(endpoint, "/")
  
    if len(parts) >= 3 && parts[1] == "api" && parts[2] == "v1" {
        endpoint = strings.Join(parts[3:], "/")
    }
    return "/" + endpoint
}


