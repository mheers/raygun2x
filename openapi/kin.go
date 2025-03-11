package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"raygun/types"

	"github.com/aws/smithy-go/ptr"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func CreateOpenAPISpec(suites []types.TestSuite) (*openapi3.T, error) {
	// Initialize the OpenAPI specification
	swagger := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Raygun API",
			Version:     "1.0.0",
			Description: "API generated from Raygun test suites",
		},
		Paths:      openapi3.NewPaths(),
		Components: &openapi3.Components{},
	}

	// Map to track paths we've already added to avoid duplicates
	pathMethods := make(map[string]map[string]bool)

	for _, suite := range suites {
		for _, test := range suite.Tests {
			path := test.DecisionPath
			method := http.MethodPost // From your code, it seems all tests use POST

			// Initialize method map if needed
			if _, exists := pathMethods[path]; !exists {
				pathMethods[path] = make(map[string]bool)
			}

			// Get the path item or create a new one
			pathItem := swagger.Paths.Find(path)
			if pathItem == nil {
				pathItem = &openapi3.PathItem{}
				swagger.Paths.Set(path, pathItem)
			}

			// Skip if we've already processed this path+method
			if pathMethods[path][method] {
				continue
			}
			pathMethods[path][method] = true

			// Create a new operation
			operation := &openapi3.Operation{
				Summary:     test.Name,
				Description: test.Description,
				Tags:        []string{suite.Name},
				Responses:   openapi3.NewResponses(),
			}

			// Parse the test input as JSON to create a more accurate request schema
			var inputObj interface{}
			requestSchemaRef := &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
				},
			}

			if err := json.Unmarshal([]byte(test.Input.Value), &inputObj); err == nil {
				// Create schema from the parsed JSON
				requestSchemaRef = createSchemaFromJSON(inputObj)
			}

			// Add request body
			operation.RequestBody = &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  openapi3.NewContent(),
				},
			}
			operation.RequestBody.Value.Content["application/json"] = &openapi3.MediaType{
				Schema: requestSchemaRef,
			}

			// Add response schema
			responseSchema := &openapi3.Schema{
				Type:       &openapi3.Types{"object"},
				Properties: make(map[string]*openapi3.SchemaRef),
			}

			// If we have expectations, try to capture them in the response
			if len(test.Expects) > 0 {
				for i, expect := range test.Expects {
					if expect.ExpectationType == "substring" && expect.Target != "" {
						// Add as a potential response property
						propName := fmt.Sprintf("expected_%d", i+1)
						responseSchema.Properties[propName] = &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{"string"},
								Example: expect.Target,
							},
						}
					}
				}
			}

			// Create successful response
			successResponse := &openapi3.Response{
				Description: ptr.String("Successful response"),
				Content:     openapi3.NewContent(),
			}
			successResponse.Content["application/json"] = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: responseSchema,
				},
			}

			// Add 200 response
			operation.Responses.Set("200", &openapi3.ResponseRef{
				Value: successResponse,
			})

			// Create error response
			errorSchema := &openapi3.Schema{
				Type: &openapi3.Types{"object"},
				Properties: map[string]*openapi3.SchemaRef{
					"error": {
						Value: &openapi3.Schema{
							Type:    &openapi3.Types{"string"},
							Example: "Error message",
						},
					},
				},
			}

			errorResponse := &openapi3.Response{
				Description: ptr.String("Error response"),
				Content:     openapi3.NewContent(),
			}
			errorResponse.Content["application/json"] = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: errorSchema,
				},
			}

			// Add 400 response
			operation.Responses.Set("400", &openapi3.ResponseRef{
				Value: errorResponse,
			})

			// Set the operation on the path item
			switch method {
			case http.MethodPost:
				pathItem.Post = operation
			case http.MethodGet:
				pathItem.Get = operation
			case http.MethodPut:
				pathItem.Put = operation
			case http.MethodDelete:
				pathItem.Delete = operation
			}

			// Update the path item in the paths
			swagger.Paths.Set(path, pathItem)
		}
	}

	// If no paths were added, add a placeholder
	if len(pathMethods) == 0 {
		placeholderPath := "/api"
		placeholderOp := &openapi3.Operation{
			Summary:     "Placeholder API",
			Description: "This is a placeholder endpoint",
			Responses:   openapi3.NewResponses(),
		}
		placeholderOp.Responses.Set("200", &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: ptr.String("Successful response"),
			},
		})

		placeholderPathItem := &openapi3.PathItem{
			Get: placeholderOp,
		}
		swagger.Paths.Set(placeholderPath, placeholderPathItem)
	}

	return swagger, nil
}

// Helper function to create a schema from a JSON object
func createSchemaFromJSON(data interface{}) *openapi3.SchemaRef {
	switch v := data.(type) {
	case map[string]interface{}:
		schema := &openapi3.Schema{
			Type:       &openapi3.Types{"object"},
			Properties: make(map[string]*openapi3.SchemaRef),
		}

		for key, val := range v {
			schema.Properties[key] = createSchemaFromJSON(val)
		}

		return &openapi3.SchemaRef{Value: schema}

	case []interface{}:
		if len(v) > 0 {
			// Use the first element to determine the array item type
			itemSchema := createSchemaFromJSON(v[0])
			return &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:  &openapi3.Types{"array"},
					Items: itemSchema,
				},
			}
		}

		// Empty array - use string as default item type
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"array"},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
				},
			},
		}

	case string:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    &openapi3.Types{"string"},
				Example: v,
			},
		}

	case float64:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    &openapi3.Types{"number"},
				Example: v,
			},
		}

	case bool:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:    &openapi3.Types{"boolean"},
				Example: v,
			},
		}

	case nil:
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"null"},
			},
		}

	default:
		// Default to string for unknown types
		return &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{"string"},
			},
		}
	}
}

func MarshalSpec(c *openapi3.T) ([]byte, error) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("error marshaling collection to YAML: %v", err)
	}

	return bytes, nil
}
