// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/entries/": {
            "get": {
                "description": "Retrieves the list of all entries from the database.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Get all entries",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Entry"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new entry. Expects a JSON object in the request body.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Create a new entry",
                "parameters": [
                    {
                        "description": "Entry information",
                        "name": "entry",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Entry"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Entry"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/entries/author": {
            "get": {
                "description": "Retrieves entries authored by the given author(s).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Get entries by author",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Author to search",
                        "name": "author",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Entry"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/entries/date": {
            "get": {
                "description": "Retrieves entries created on the given date.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Get entries by date",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Creation date (YYYY-MM-DD)",
                        "name": "createdAt",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Entry"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid date format. Expected YYYY-MM-DD",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/entries/health": {
            "get": {
                "description": "Checks if the service is up",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/entries/id/": {
            "get": {
                "description": "Retrieves an entry by its ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Get an entry by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Entry ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Entry"
                        }
                    },
                    "400": {
                        "description": "Invalid ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Entry not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an entry by its ID. Expects a JSON object in the request body.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Update an entry by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Entry ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Updated entry information",
                        "name": "entry",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Entry"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Entry"
                        }
                    },
                    "400": {
                        "description": "Invalid ID or request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Entry not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes an entry by its ID.",
                "tags": [
                    "Entries"
                ],
                "summary": "Delete an entry by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Entry ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Entry not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/entries/title": {
            "get": {
                "description": "Retrieves entries that match the given title.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Get entries by title",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Title to search",
                        "name": "title",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Entry"
                        }
                    },
                    "404": {
                        "description": "Entry not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/entries/wiki": {
            "get": {
                "description": "Retrieves entries associated with a specific Wiki ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entries"
                ],
                "summary": "Get entries by Wiki ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Wiki ID",
                        "name": "wikiID",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Entry"
                            }
                        }
                    },
                    "400": {
                        "description": "WikiID is required",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Entry": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "wiki_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8002",
	BasePath:         "/api/entries",
	Schemes:          []string{},
	Title:            "Entry Service API",
	Description:      "API documentation for the Entry Service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
