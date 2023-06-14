// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/example/helloworld": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "ping example",
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
        "/login": {
            "get": {
                "description": "This route handles the '/login' endpoint and initiates OAuth authentication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "Initiate OAuth authentication",
                "parameters": [
                    {
                        "type": "string",
                        "description": "OAuth provider name",
                        "name": "provider",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirection to the OAuth URL",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/oauth-redirect": {
            "get": {
                "description": "This route handles the '/oauth-redirect' endpoint and finalizes the OAuth authentication process. After successful authentication, the session is updated with a token and a userId.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "Finalize OAuth authentication",
                "parameters": [
                    {
                        "type": "string",
                        "description": "OAuth code received from provider",
                        "name": "code",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "OAuth state received from provider",
                        "name": "state",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/v1/tmdb/movies": {
            "get": {
                "description": "Get movies that match the specified genre and were released in a specific year",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get movies by genre and release date",
                "operationId": "get-movies-by-genre-date",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Genre ID",
                        "name": "genre",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Release Year",
                        "name": "year",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Movie"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/v1/tmdb/movies/post/ids": {
            "post": {
                "description": "Get movies with given IDs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get movies by IDs",
                "operationId": "get-movies-by-ids",
                "parameters": [
                    {
                        "description": "List of Movie IDs",
                        "name": "ids",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Movie"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/v1/tmdb/movies/post/title": {
            "post": {
                "description": "Get movies that match given titles",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Search movies by title",
                "operationId": "get-movies-by-title",
                "parameters": [
                    {
                        "description": "List of Titles",
                        "name": "titles",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Movie"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        },
        "/v1/tmdb/random10": {
            "get": {
                "description": "Get 10 random popular movies",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get random popular movies",
                "operationId": "get-random-movies",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Movie"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "types.Genre": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "types.Movie": {
            "type": "object",
            "properties": {
                "adult": {
                    "type": "boolean"
                },
                "backdrop_path": {
                    "type": "string"
                },
                "belongs_to_collection": {},
                "budget": {
                    "type": "integer"
                },
                "director": {
                    "type": "string"
                },
                "genres": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Genre"
                    }
                },
                "homepage": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "imdb_id": {
                    "type": "string"
                },
                "language": {
                    "type": "string"
                },
                "original_language": {
                    "type": "string"
                },
                "original_title": {
                    "type": "string"
                },
                "overview": {
                    "type": "string"
                },
                "popularity": {
                    "type": "number"
                },
                "poster_path": {
                    "type": "string"
                },
                "production_companies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.ProductionCompany"
                    }
                },
                "production_countries": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.ProductionCountry"
                    }
                },
                "release_date": {
                    "type": "string"
                },
                "revenue": {
                    "type": "integer"
                },
                "runtime": {
                    "type": "integer"
                },
                "spoken_languages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.SpokenLanguage"
                    }
                },
                "status": {
                    "type": "string"
                },
                "tagline": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "video": {
                    "type": "boolean"
                },
                "vote_average": {
                    "type": "number"
                },
                "vote_count": {
                    "type": "integer"
                }
            }
        },
        "types.ProductionCompany": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "logo_path": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "origin_country": {
                    "type": "string"
                }
            }
        },
        "types.ProductionCountry": {
            "type": "object",
            "properties": {
                "iso_3166_1": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "types.SpokenLanguage": {
            "type": "object",
            "properties": {
                "english_name": {
                    "type": "string"
                },
                "iso_639_1": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
