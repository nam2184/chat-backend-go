package Errorors

import "net/http"

type ErrorType struct {
	Section string `json:"section"`
	Message string `json:"message"`
	err   error  
}

func NewErrorType(req *http.Request, err error) ErrorType {
  ctx := req.Context()    
  section, ok := ctx.Value("section").(string); if !ok {
    return ErrorType{ Message: err.Error(), err: err} 
  }
  return ErrorType{Section : section, Message: err.Error(), err : err,}
}


func (e ErrorType) NAME() string {
  return "error_type"
}
