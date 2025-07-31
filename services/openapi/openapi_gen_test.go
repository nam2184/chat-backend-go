package openapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/nam2184/mymy/middleware"
	"github.com/nam2184/mymy/models/response"
)

func TestGenerateWholeSchema(t *testing.T) {
	components := &openapi3.Components{Schemas: openapi3.Schemas{}}

	skipParam := &openapi3.ParameterRef{Value: openapi3.NewQueryParameter("skip").WithSchema(openapi3.NewIntegerSchema()).WithRequired(false)}
	limitParam := &openapi3.ParameterRef{Value: openapi3.NewQueryParameter("limit").WithSchema(openapi3.NewIntegerSchema()).WithRequired(false)}

	addSchema := func(name string, model interface{}) {
		schemaRef, err := openapi3gen.NewSchemaRefForValue(model, nil)
		if err != nil {
			t.Fatalf("failed to generate schema for %s: %v", name, err)
		}
		components.Schemas[name] = schemaRef
	}
	addSchema("GetMessages", response.StandardResponse[response.GetMessages]{})
	addSchema("GetEncryptedMessages", response.StandardResponse[response.GetEncryptedMessages]{})
	addSchema("GetChats", response.GetChats{})
	addSchema("GetUser", response.GetUser{})
	addSchema("PostAuthResponse", response.PostAuthResponse{})
	addSchema("PostSignUpResponse", response.PostSignUpResponse{})
	addSchema("GetRefreshAuthResponse", response.GetRefreshAuthResponse{})
	addSchema("ErrorResponse", middleware.ErrorResponse{})

	createResponseWithSchema := func(description string, schemaRef *openapi3.SchemaRef, errorRef *openapi3.SchemaRef) *openapi3.Responses {
		responses := openapi3.NewResponses()
		responses.Set("200", &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: strPtr(description),
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: schemaRef,
					},
				},
			},
		})
		responses.Set("400", &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: strPtr("Failure"),
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: errorRef,
					},
				},
			},
		})
		return responses
	}

	paths := openapi3.NewPaths(
		openapi3.WithPath("/auth", &openapi3.PathItem{
			Post: &openapi3.Operation{
				Summary:     "Authenticate user and return tokens",
				OperationID: "PostAuth",
				Responses:   createResponseWithSchema("Authentication success", &openapi3.SchemaRef{Ref: "#/components/schemas/PostAuthResponse"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),
		openapi3.WithPath("/auth/signup", &openapi3.PathItem{
			Post: &openapi3.Operation{
				Summary:     "Create new account",
				OperationID: "PostAuthSignUp",
				Responses:   createResponseWithSchema("Create account success", &openapi3.SchemaRef{Ref: "#/components/schemas/PostSignUpResponse"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),
		openapi3.WithPath("/auth/refresh", &openapi3.PathItem{
			Get: &openapi3.Operation{
				Summary:     "Refresh access token",
				OperationID: "GetRefreshToken",
				Responses:   createResponseWithSchema("Refresh Authentication success", &openapi3.SchemaRef{Ref: "#/components/schemas/GetRefreshAuthResponse"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),

		openapi3.WithPath("/user", &openapi3.PathItem{
			Get: &openapi3.Operation{
				Summary:     "Get current user info",
				OperationID: "GetUser",
				Responses:   createResponseWithSchema("User retrieved", &openapi3.SchemaRef{Ref: "#/components/schemas/GetUser"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),

		openapi3.WithPath("/chats", &openapi3.PathItem{
			Get: &openapi3.Operation{
				Summary:     "Get chat list",
				OperationID: "GetChats",
				Responses:   createResponseWithSchema("Chats retrieved", &openapi3.SchemaRef{Ref: "#/components/schemas/GetChats"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),

		openapi3.WithPath("/messages/{chatID}", &openapi3.PathItem{
			Get: &openapi3.Operation{
				Summary:     "Get messages for a chat",
				OperationID: "GetMessages",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: &openapi3.Parameter{
							Name:     "chatID",
							In:       "path",
							Required: true,
							Schema:   openapi3.NewIntegerSchema().NewRef(),
						},
					}, skipParam, limitParam,
				},
				Responses: createResponseWithSchema("Messages retrieved", &openapi3.SchemaRef{Ref: "#/components/schemas/GetMessages"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),
		openapi3.WithPath("/encrypted-messages/{chatID}", &openapi3.PathItem{
			Get: &openapi3.Operation{
				Summary:     "Get messages for a chat",
				OperationID: "GetEncryptedMessages",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: &openapi3.Parameter{
							Name:     "chatID",
							In:       "path",
							Required: true,
							Schema:   openapi3.NewIntegerSchema().NewRef(),
						},
					}, skipParam, limitParam,
				},
				Responses: createResponseWithSchema("Encrypted Messages retrieved", &openapi3.SchemaRef{Ref: "#/components/schemas/GetEncryptedMessages"}, &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),

		openapi3.WithPath("/messages/{chatID}/count", &openapi3.PathItem{
			Get: &openapi3.Operation{
				Summary:     "Get message count",
				OperationID: "GetMessagesCount",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: &openapi3.Parameter{
							Name:     "chatID",
							In:       "path",
							Required: true,
							Schema:   openapi3.NewIntegerSchema().NewRef(),
						},
					}, skipParam, limitParam,
				},
				Responses: createResponseWithSchema("Count retrieved", openapi3.NewIntegerSchema().NewRef(), &openapi3.SchemaRef{Ref: "#/components/schemas/ErrorResponse"}),
			},
		}),
	)

	servers := openapi3.Servers{
		{URL: "https://hello.example.com", Description: " Server"},
		{URL: "http://localhost:8000/api/v1", Description: "Local Test Server"},
	}

	doc := &openapi3.T{
		OpenAPI:    "3.0.0",
		Info:       &openapi3.Info{Title: "Dem API", Version: "0.1"},
		Servers:    servers,
		Paths:      paths,
		Components: components,
	}

	docj, _ := doc.MarshalJSON()
	var formattedJSON bytes.Buffer
	_ = json.Indent(&formattedJSON, docj, "", "  ")
	_ = ioutil.WriteFile("../../schema.json", formattedJSON.Bytes(), 0644)
}

func strPtr(s string) *string {
	return &s
}
