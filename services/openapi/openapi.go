package openapi 

import (
  "github.com/getkin/kin-openapi/openapi3"
  "os"
)

func LinkSchema(test bool) (*openapi3.T, *openapi3.Loader, error) {
  var file string
  if test != true {
    file = "/test_schema.json"
  } else {
    file = "/test_schema2.json"
  }
  dir, err := os.Getwd()

  openapi, err := os.ReadFile(dir + file)
  if err != nil {
    return nil, nil, err
  }
  // Parse the OpenAPI document
  loader := openapi3.NewLoader()
  doc, err := loader.LoadFromData(openapi)
  if err != nil {
    return nil, nil, err 
  }
  // Validate the OpenAPI document
  if err := doc.Validate(loader.Context); err != nil {
    return nil, nil, err
  }

  return doc, loader, nil
}
