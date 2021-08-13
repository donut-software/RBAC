package rest

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

//go:generate go run ../../cmd/openapi-gen/main.go -path .
// //go:generate oapi-codegen -package openapi3 -generate types  -o ../../pkg/openapi3/rbac_types.gen.go openapi3.yaml
// //go:generate oapi-codegen -package openapi3 -generate client -o ../../pkg/openapi3/client.gen.go     openapi3.yaml

// NewOpenAPI3 instantiates the OpenAPI specification for this service.
func NewOpenAPI3() openapi3.Swagger {
	swagger := openapi3.Swagger{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "RBAC API",
			Description: "REST APIs used for interacting with the RBAC Service",
			Version:     "0.0.0",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://192.168.10.199:9234",
			},
		},
	}

	swagger.Components.Schemas = openapi3.Schemas{
		"Profile": openapi3.NewSchemaRef("",
			openapi3.NewObjectSchema().
				WithProperty("id", openapi3.NewUUIDSchema()).
				WithProperty("profile_picture", openapi3.NewStringSchema()).
				WithProperty("profile_background", openapi3.NewStringSchema()).
				WithProperty("first_name", openapi3.NewStringSchema()).
				WithProperty("last_name", openapi3.NewStringSchema()).
				WithProperty("mobile", openapi3.NewStringSchema()).
				WithProperty("email", openapi3.NewStringSchema()).
				WithProperty("is_blocked", openapi3.NewBoolSchema()).
				WithProperty("created_at", openapi3.NewStringSchema().WithFormat("date-time"))),
		"Account": openapi3.NewSchemaRef("",
			openapi3.NewObjectSchema().
				WithProperty("id", openapi3.NewUUIDSchema()).
				WithProperty("username", openapi3.NewStringSchema()).
				WithProperty("created_at", openapi3.NewStringSchema().WithFormat("date-time")).
				WithPropertyRef("profile", &openapi3.SchemaRef{
					Ref: "#/components/schemas/Profile",
				})),
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"CreateAccountRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for registering an account.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("username", openapi3.NewStringSchema()).
					WithProperty("password", openapi3.NewStringSchema()).
					WithProperty("profile_picture", openapi3.NewStringSchema()).
					WithProperty("profile_background", openapi3.NewStringSchema()).
					WithProperty("first_name", openapi3.NewStringSchema()).
					WithProperty("last_name", openapi3.NewStringSchema()).
					WithProperty("mobile", openapi3.NewStringSchema()).
					WithProperty("email", openapi3.NewStringSchema())),
		},
		"GetAccountRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for registering an account.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("username", openapi3.NewStringSchema()).
					WithProperty("password", openapi3.NewStringSchema()).
					WithProperty("profile_picture", openapi3.NewStringSchema()).
					WithProperty("profile_background", openapi3.NewStringSchema()).
					WithProperty("first_name", openapi3.NewStringSchema()).
					WithProperty("last_name", openapi3.NewStringSchema()).
					WithProperty("mobile", openapi3.NewStringSchema()).
					WithProperty("email", openapi3.NewStringSchema())),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response when errors happen.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("error", openapi3.NewStringSchema()))),
		},
		"CreateAccountResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned back after registering an accounts.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("message", openapi3.NewStringSchema()))),
		},
		"GetAccountResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned back after registering an accounts.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("account", &openapi3.SchemaRef{
						Ref: "#/components/schemas/Account",
					}))),
		},
		// "ReadTasksResponse": &openapi3.ResponseRef{
		// 	Value: openapi3.NewResponse().
		// 		WithDescription("Response returned back after searching one task.").
		// 		WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
		// 			WithProperty("message", openapi3.NewStringSchema()))),
		// },
	}

	swagger.Paths = openapi3.Paths{
		"/register": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "RegisterAccount",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/CreateAccountRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/CreateAccountResponse",
					},
				},
			},
		},
		"/accounts/{username}": &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "Get Account",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: openapi3.NewPathParameter("username").
							WithSchema(openapi3.NewStringSchema()),
					},
				},
				Responses: openapi3.Responses{
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/GetAccountResponse",
					},
				},
			},
		},
		// 	Put: &openapi3.Operation{
		// 		OperationID: "UpdateTask",
		// 		Parameters: []*openapi3.ParameterRef{
		// 			{
		// 				Value: openapi3.NewPathParameter("taskId").
		// 					WithSchema(openapi3.NewUUIDSchema()),
		// 			},
		// 		},
		// 		RequestBody: &openapi3.RequestBodyRef{
		// 			Ref: "#/components/requestBodies/UpdateTasksRequest",
		// 		},
		// 		Responses: openapi3.Responses{
		// 			"400": &openapi3.ResponseRef{
		// 				Ref: "#/components/responses/ErrorResponse",
		// 			},
		// 			"500": &openapi3.ResponseRef{
		// 				Ref: "#/components/responses/ErrorResponse",
		// 			},
		// 			"200": &openapi3.ResponseRef{
		// 				Value: openapi3.NewResponse().WithDescription("Task was updated"),
		// 			},
		// 		},
		// 	},
		// },
	}

	return swagger
}

func RegisterOpenAPI(r *mux.Router) {
	swagger := NewOpenAPI3()

	r.HandleFunc("/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		renderResponse(w, &swagger, http.StatusOK)
	}).Methods(http.MethodGet)

	r.HandleFunc("/openapi3.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")

		data, _ := yaml.Marshal(&swagger)

		_, _ = w.Write(data)

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
}
