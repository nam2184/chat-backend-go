package openapi

import (
	"bytes"
	//qm "github.com/nam2184/mymy/models/db"
	//rs "github.com/nam2184/mymy/models/response"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGenerateWholeSchema(t *testing.T) {
	//models_db := []interface{}{
	//qm.Auth{},
	//qm.Chat{},
	//qm.User{},
	//qm.Message{},
	//}

	servers := openapi3.Servers{
		{
			URL:         "http://api.app.internal/api/v1",
			Description: "Test Server",
		},
		{
			URL:         "http://localhost:8000/api/v1",
			Description: "Local Test Server",
		},
	}
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Dem API",
			Version: "0.1",
		},
		Servers: servers,
		Paths: openapi3.NewPaths(
			openapi3.WithPath("/auth", &openapi3.PathItem{}),
			openapi3.WithPath("/auth/refresh", &openapi3.PathItem{}),
			openapi3.WithPath("/messages", &openapi3.PathItem{}),
			openapi3.WithPath("/ws", &openapi3.PathItem{}),
		),
	}

	docj, _ := doc.MarshalJSON()
	var formattedJSON bytes.Buffer
	_ = json.Indent(&formattedJSON, docj, "", "  ")
	fileName := "test_schema.json"
	_ = ioutil.WriteFile("../../"+fileName, formattedJSON.Bytes(), 0644)
	_ = ioutil.WriteFile("../../routes/"+fileName, formattedJSON.Bytes(), 0644)

}
