// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "r4ha",
            "url": "https://github.com/rahadinabudiman"
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
        "/admin/user": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin - User"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search data",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "asc",
                            "desc"
                        ],
                        "type": "string",
                        "description": "Sort by name",
                        "name": "sortBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GetAllUserStatusOKResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.BadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.UnauthorizedResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/dtos.ForbiddenResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dtos.NotFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.InternalServerErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Login an account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User - Auth"
                ],
                "summary": "Login User with Username and Password",
                "parameters": [
                    {
                        "description": "Payload Body [RAW]",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.LoginStatusOKResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.BadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.UnauthorizedResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/dtos.ForbiddenResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dtos.NotFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.InternalServerErrorResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Login an account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User - Auth"
                ],
                "summary": "Login User with Username and Password",
                "parameters": [
                    {
                        "description": "Payload Body [RAW]",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.UserRegister"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.RegisterStatusOKResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.BadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.UnauthorizedResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/dtos.ForbiddenResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dtos.NotFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.InternalServerErrorResponse"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "User get profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User - Account"
                ],
                "summary": "Get Profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ProfileStatusOKResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.BadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.UnauthorizedResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/dtos.ForbiddenResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dtos.NotFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.InternalServerErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "User update an Profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User - Account"
                ],
                "summary": "Update Profile",
                "parameters": [
                    {
                        "description": "Payload Body [RAW]",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.UserUpdateProfile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ProfileStatusOKResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.BadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.UnauthorizedResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/dtos.ForbiddenResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dtos.NotFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.InternalServerErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "User delete an Profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User - Account"
                ],
                "summary": "Delete Profile",
                "parameters": [
                    {
                        "description": "Payload Body [RAW]",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.DeleteUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.StatusOKResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.BadRequestResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.UnauthorizedResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/dtos.ForbiddenResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dtos.NotFoundResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dtos.InternalServerErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.BadRequestResponse": {
            "type": "object",
            "properties": {
                "errors": {},
                "message": {
                    "type": "string",
                    "example": "Bad Request"
                },
                "status_code": {
                    "type": "integer",
                    "example": 400
                }
            }
        },
        "dtos.DeleteUserRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 6,
                    "example": "rahadinabudimansundara"
                }
            }
        },
        "dtos.ForbiddenResponse": {
            "type": "object",
            "properties": {
                "errors": {},
                "message": {
                    "type": "string",
                    "example": "Forbidden"
                },
                "status_code": {
                    "type": "integer",
                    "example": 403
                }
            }
        },
        "dtos.GetAllUserStatusOKResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dtos.UserDetailResponse"
                },
                "message": {
                    "type": "string",
                    "example": "Successfully get profile"
                },
                "meta": {
                    "$ref": "#/definitions/utils.Meta"
                },
                "status_code": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "dtos.InternalServerErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {},
                "message": {
                    "type": "string",
                    "example": "Internal Server Error"
                },
                "status_code": {
                    "type": "integer",
                    "example": 500
                }
            }
        },
        "dtos.LoginStatusOKResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dtos.UserLoginResponse"
                },
                "message": {
                    "type": "string",
                    "example": "Login Success"
                },
                "status_code": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "dtos.NotFoundResponse": {
            "type": "object",
            "properties": {
                "errors": {},
                "message": {
                    "type": "string",
                    "example": "Not Found"
                },
                "status_code": {
                    "type": "integer",
                    "example": 404
                }
            }
        },
        "dtos.ProfileStatusOKResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dtos.UserProfileResponse"
                },
                "message": {
                    "type": "string",
                    "example": "Get Profile Success"
                },
                "status_code": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "dtos.RegisterStatusOKResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dtos.UserRegisterResponse"
                },
                "message": {
                    "type": "string",
                    "example": "Register Success"
                },
                "status_code": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "dtos.StatusOKResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string",
                    "example": "Successfully"
                },
                "status_code": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "dtos.UnauthorizedResponse": {
            "type": "object",
            "properties": {
                "errors": {},
                "message": {
                    "type": "string",
                    "example": "Unauthorized"
                },
                "status_code": {
                    "type": "integer",
                    "example": 401
                }
            }
        },
        "dtos.UserDetailResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dtos.UserLoginRequest": {
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
        "dtos.UserLoginResponse": {
            "type": "object",
            "required": [
                "token",
                "username"
            ],
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
                },
                "username": {
                    "type": "string",
                    "example": "rahadinabudimansundara"
                }
            }
        },
        "dtos.UserProfileResponse": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dtos.UserRegister": {
            "type": "object",
            "required": [
                "confirm_password",
                "name",
                "password",
                "username"
            ],
            "properties": {
                "confirm_password": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "example": "Admin"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dtos.UserRegisterResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dtos.UserUpdateProfile": {
            "type": "object",
            "required": [
                "name",
                "username"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.Meta": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "integer",
                    "example": 1
                },
                "next_page": {},
                "prev_page": {
                    "type": "integer",
                    "example": 1
                },
                "total": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "54.179.176.114:1309",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "FiberStore Documentation API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}