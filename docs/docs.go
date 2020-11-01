// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Brian Campos Castro",
            "email": "brian.campos.castro@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "description": "returns the time from DB",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "models": [
                    "commons"
                ],
                "summary": "returns the time from DB",
                "parameters": [
                    {
                        "type": "string",
                        "description": "API_KEY Header",
                        "name": "API_KEY",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/models": {
            "get": {
                "description": "Get all the models by defining the fields",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "models": [
                    "models"
                ],
                "summary": "Get all models",
                "parameters": [
                    {
                        "type": "string",
                        "description": "API_KEY Header",
                        "name": "API_KEY",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "Event fields",
                        "name": "fields",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new model item",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "models": [
                    "models"
                ],
                "summary": "Create a model",
                "parameters": [
                    {
                        "description": "New model",
                        "name": "model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/validators.CreateBody"
                        }
                    },
                    {
                        "type": "string",
                        "description": "API_KEY Header",
                        "name": "API_KEY",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/models/list": {
            "get": {
                "description": "Get all the models by ids",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "models": [
                    "models"
                ],
                "summary": "Get all models by ids",
                "parameters": [
                    {
                        "type": "string",
                        "description": "API_KEY Header",
                        "name": "API_KEY",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "Models IDs",
                        "name": "modelIds",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/models/{id}": {
            "get": {
                "description": "Get a model by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "models": [
                    "models"
                ],
                "summary": "Get a model",
                "parameters": [
                    {
                        "type": "string",
                        "description": "API_KEY Header",
                        "name": "API_KEY",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Model ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a model item",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "models": [
                    "models"
                ],
                "summary": "Update a model",
                "parameters": [
                    {
                        "type": "string",
                        "description": "API_KEY Header",
                        "name": "API_KEY",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "model ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "model Updated",
                        "name": "model",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/validators.UpdateBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "utils.ResponseData": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "href": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "utils.ResponseMessage": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/utils.ResponseData"
                    }
                },
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/utils.ResponseData"
                    }
                }
            }
        },
        "validators.CreateBody": {
            "type": "object",
            "properties": {
                "insUserId": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "validators.UpdateBody": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "setUserId": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:3040",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "4g Emailing API",
	Description: "Games Provider API is a 4G Solution. Responsible of managing information emails..",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
