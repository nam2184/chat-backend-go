package response

import (
	qm "github.com/nam2184/generic-queries"
	e "github.com/nam2184/mymy/models/errors"
)

func StandardSuccessResponse[T any](meta Meta, array ...T) *StandardResponse[T] {
  return &StandardResponse[T]{ Array : array, Meta: meta }
}

func QuerySuccessResponse[T qm.QueryTypes](message string, payload []T) *QueryResponse[T] {
  return &QueryResponse[T]{ Status: "success", Type: "content", Message: message, Payload: payload }
}

func QueryFailureResponse[T qm.QueryTypes] (message string, error e.ErrorType) *QueryResponse[T] {
  return &QueryResponse[T]{ Status: "failure", Type: "content", Error: error}
}




