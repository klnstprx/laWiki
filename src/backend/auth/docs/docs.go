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
        "/api/auth/callback": {
            "get": {
                "description": "Handles the OAuth2 callback from Google",
                "tags": [
                    "Authentication"
                ],
                "summary": "OAuth2 Callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "OAuth state",
                        "name": "state",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authorization code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirect after login",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Invalid OAuth state",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Could not create JWT",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/health": {
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
        "/api/auth/login": {
            "get": {
                "description": "Initiates the OAuth2 flow with Google",
                "tags": [
                    "Authentication"
                ],
                "summary": "Initiate OAuth2 Login",
                "responses": {
                    "302": {
                        "description": "Redirect to Google OAuth2 login",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "OAuth2": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://accounts.google.com/o/oauth2/auth",
            "tokenUrl": "https://oauth2.googleapis.com/token"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/auth",
	Schemes:          []string{},
	Title:            "Auth Service API",
	Description:      "API documentation for the Auth Service. !!THIS SERVICE IS WIP!!",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
