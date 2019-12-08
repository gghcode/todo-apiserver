// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-12-08 04:57:52.70247027 +0000 UTC m=+0.087265640

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
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/refresh": {
            "post": {
                "description": "Get new access token by refreshtoken",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth API"
                ],
                "parameters": [
                    {
                        "description": "payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/auth.AccessTokenByRefreshRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/auth.tokenResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credential",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/token": {
            "post": {
                "description": "Issue new token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth API"
                ],
                "parameters": [
                    {
                        "description": "payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/auth.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credential",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/healthy": {
            "get": {
                "description": "Get server healthy",
                "tags": [
                    "App API"
                ],
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
        "/api/version": {
            "get": {
                "description": "Get server version",
                "tags": [
                    "App API"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorMetaData": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "$ref": "#/definitions/api.ErrorMetaData"
                }
            }
        },
        "auth.AccessTokenByRefreshRequest": {
            "type": "object",
            "required": [
                "token"
            ],
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "auth.TokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "auth.tokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "APAS TODO API",
	Description: "This is a apas todo api server.",
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
