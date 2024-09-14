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
        "/admin/category": {
            "get": {
                "description": "Retrieves a list of all categories",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Get all categories",
                "responses": {}
            },
            "post": {
                "description": "Allows admin to add a new category",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Add a new category",
                "parameters": [
                    {
                        "description": "Category Name",
                        "name": "category_name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/helper.AddCategory"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/admin/category/{CatID}": {
            "delete": {
                "description": "Deletes a category by ID, used by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Admin delete skill",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Category ID",
                        "name": "CatID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/admin/login": {
            "post": {
                "description": "Log in as an admin using email and password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Admin login",
                "parameters": [
                    {
                        "description": "Admin login credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/helper.ADLogin"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/admin/skill": {
            "post": {
                "description": "This endpoint allows a admin to add a new skill by providing the skill name.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Add a new skill to a user's profile",
                "parameters": [
                    {
                        "description": "Skill information",
                        "name": "skill",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/helper.AddSkill"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/admin/skills": {
            "get": {
                "description": "Retrieves a list of all skills",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Get all skills",
                "responses": {}
            }
        },
        "/admin/skills/{skillID}": {
            "delete": {
                "description": "Deletes a skill by ID, used by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Admin delete skill",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Skill ID",
                        "name": "skillID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/admin/users": {
            "get": {
                "description": "Retrieve a list of all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Get all users",
                "responses": {}
            }
        },
        "/gig/add": {
            "post": {
                "description": "This endpoint creates a new gig with a title, description, price, and images. Images are uploaded via multipart form.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gigs"
                ],
                "summary": "Create a new gig",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Title of the gig",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Description of the gig",
                        "name": "description",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Category of the gig",
                        "name": "category",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Number of delivery days",
                        "name": "delivery",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Number of revisions",
                        "name": "revisions",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Price of the gig",
                        "name": "price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "file"
                        },
                        "collectionFormat": "csv",
                        "description": "Images for the gig (can upload multiple images)",
                        "name": "images",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/gig/user": {
            "get": {
                "description": "Get all gigs created by the logged-in user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Gigs"
                ],
                "summary": "Get Gigs by User ID",
                "responses": {}
            }
        },
        "/ping": {
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
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "Authenticate user with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/helper.LoginData"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/user/profile": {
            "get": {
                "description": "Retrieves the profile details of the user based on their user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get user profile",
                "responses": {}
            },
            "put": {
                "description": "Update the user's bio and title in their profile",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update user profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User bio",
                        "name": "Bio",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User title",
                        "name": "Title",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/user/profile-photo": {
            "post": {
                "description": "Uploads a profile photo for the user based on the userID.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Uploads a profile photo for the user",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Profile photo",
                        "name": "photo",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/user/signup": {
            "post": {
                "description": "Create a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "Signup Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/helper.SignupData"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/user/skill/{Skill}": {
            "delete": {
                "description": "Deletes a specific skill for a user based on the user ID and skill ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Delete a skill from a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Skill ID to delete",
                        "name": "Skill",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/user/skills": {
            "post": {
                "description": "This endpoint allows a freelancer to add a skill and set their proficiency level.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Add freelancer skill",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Skill name",
                        "name": "skillName",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Proficiency level",
                        "name": "proficency",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "helper.ADLogin": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "helper.AddCategory": {
            "type": "object",
            "properties": {
                "category_name": {
                    "type": "string"
                }
            }
        },
        "helper.AddSkill": {
            "type": "object",
            "properties": {
                "skill_name": {
                    "type": "string"
                }
            }
        },
        "helper.LoginData": {
            "type": "object",
            "properties": {
                "useremail": {
                    "type": "string"
                },
                "userpassword": {
                    "type": "string"
                }
            }
        },
        "helper.SignupData": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string"
                },
                "firstname": {
                    "type": "string"
                },
                "lastname": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "useremail": {
                    "type": "string"
                },
                "userpassword": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "API Gateway Swagger",
	Description:      "This is the API Gateway for the Flexi Worke project",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
