package docs

import "github.com/swaggo/swag"

// @title Vite Plugin End API
// @version 1.0
// @description Vite插件管理系统的API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请在值前加上"Bearer "前缀，例如："Bearer abcde12345"

// @tag.name 用户管理
// @tag.description 用户注册、登录、信息管理等接口

// @tag.name 插件管理
// @tag.description 插件的创建、查询、打包等接口

// @tag.name 文件管理
// @tag.description 文件上传、下载、删除等接口

// @tag.name 系统管理
// @tag.description 系统状态、健康检查等接口

func SwaggerInfo() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate: docTemplate,
	})
}

const docTemplate = `{
    "schemes": ["http"],
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/auth/register": {
            "post": {
                "tags": ["用户管理"],
                "summary": "用户注册",
                "description": "注册新用户",
                "operationId": "register",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "description": "用户注册信息",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "注册成功",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "tags": ["用户管理"],
                "summary": "用户登录",
                "description": "用户登录并获取token",
                "operationId": "login",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "description": "用户登录信息",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功",
                        "schema": {
                            "$ref": "#/definitions/LoginResponse"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "认证失败",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/profile": {
            "get": {
                "tags": ["用户管理"],
                "summary": "获取用户信息",
                "description": "获取当前登录用户的信息",
                "operationId": "getProfile",
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功获取用户信息",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "tags": ["用户管理"],
                "summary": "更新用户信息",
                "description": "更新当前登录用户的信息",
                "operationId": "updateProfile",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "description": "用户更新信息",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "更新成功",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/plugins": {
            "post": {
                "tags": ["插件管理"],
                "summary": "创建插件",
                "description": "创建新的插件",
                "operationId": "createPlugin",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "body",
                        "name": "body",
                        "description": "插件信息",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreatePluginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "创建成功",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "get": {
                "tags": ["插件管理"],
                "summary": "获取插件列表",
                "description": "获取所有插件列表",
                "operationId": "listPlugins",
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "query",
                        "name": "page",
                        "description": "页码",
                        "required": false,
                        "type": "integer",
                        "default": 1
                    },
                    {
                        "in": "query",
                        "name": "page_size",
                        "description": "每页数量",
                        "required": false,
                        "type": "integer",
                        "default": 10
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功获取插件列表",
                        "schema": {
                            "$ref": "#/definitions/PluginListResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/plugins/{key}": {
            "get": {
                "tags": ["插件管理"],
                "summary": "获取插件详情",
                "description": "根据插件key获取插件详情",
                "operationId": "getPlugin",
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "key",
                        "description": "插件key",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功获取插件详情",
                        "schema": {
                            "$ref": "#/definitions/PluginResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "插件不存在",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "tags": ["插件管理"],
                "summary": "更新插件",
                "description": "更新插件信息",
                "operationId": "updatePlugin",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "key",
                        "description": "插件key",
                        "required": true,
                        "type": "string"
                    },
                    {
                        "in": "body",
                        "name": "body",
                        "description": "插件更新信息",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdatePluginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "更新成功",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "插件不存在",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "tags": ["插件管理"],
                "summary": "删除插件",
                "description": "删除指定插件",
                "operationId": "deletePlugin",
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "key",
                        "description": "插件key",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "删除成功",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "插件不存在",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/plugins/{key}/pack": {
            "get": {
                "tags": ["插件管理"],
                "summary": "打包插件",
                "description": "将插件打包成zip文件",
                "operationId": "packPlugin",
                "produces": ["application/zip"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "key",
                        "description": "插件key",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "打包成功",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "插件不存在",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "tags": ["文件管理"],
                "summary": "上传文件",
                "description": "上传文件到服务器",
                "operationId": "uploadFile",
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "formData",
                        "name": "file",
                        "description": "要上传的文件",
                        "required": true,
                        "type": "file"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "上传成功",
                        "schema": {
                            "$ref": "#/definitions/UploadResponse"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/upload/{filename}": {
            "get": {
                "tags": ["文件管理"],
                "summary": "获取文件",
                "description": "获取上传的文件",
                "operationId": "getFile",
                "produces": ["application/octet-stream"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "filename",
                        "description": "文件名",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功获取文件",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "文件不存在",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "tags": ["文件管理"],
                "summary": "删除文件",
                "description": "删除上传的文件",
                "operationId": "deleteFile",
                "produces": ["application/json"],
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "parameters": [
                    {
                        "in": "path",
                        "name": "filename",
                        "description": "文件名",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "删除成功",
                        "schema": {
                            "$ref": "#/definitions/Response"
                        }
                    },
                    "401": {
                        "description": "未认证",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "文件不存在",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "tags": ["系统管理"],
                "summary": "健康检查",
                "description": "检查系统健康状态",
                "operationId": "healthCheck",
                "produces": ["application/json"],
                "responses": {
                    "200": {
                        "description": "系统正常",
                        "schema": {
                            "$ref": "#/definitions/HealthResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "RegisterRequest": {
            "type": "object",
            "required": ["username", "password", "email"],
            "properties": {
                "username": {
                    "type": "string",
                    "description": "用户名"
                },
                "password": {
                    "type": "string",
                    "description": "密码"
                },
                "email": {
                    "type": "string",
                    "description": "邮箱"
                }
            }
        },
        "LoginRequest": {
            "type": "object",
            "required": ["username", "password"],
            "properties": {
                "username": {
                    "type": "string",
                    "description": "用户名"
                },
                "password": {
                    "type": "string",
                    "description": "密码"
                }
            }
        },
        "LoginResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "description": "状态码"
                },
                "message": {
                    "type": "string",
                    "description": "状态信息"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "token": {
                            "type": "string",
                            "description": "JWT令牌"
                        },
                        "user": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            }
        },
        "UserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "description": "用户ID"
                },
                "username": {
                    "type": "string",
                    "description": "用户名"
                },
                "email": {
                    "type": "string",
                    "description": "邮箱"
                },
                "role": {
                    "type": "string",
                    "description": "角色"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time",
                    "description": "创建时间"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time",
                    "description": "更新时间"
                }
            }
        },
        "UpdateUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "description": "邮箱"
                }
            }
        },
        "CreatePluginRequest": {
            "type": "object",
            "required": ["name", "version", "description"],
            "properties": {
                "name": {
                    "type": "string",
                    "description": "插件名称"
                },
                "version": {
                    "type": "string",
                    "description": "插件版本"
                },
                "description": {
                    "type": "string",
                    "description": "插件描述"
                },
                "author": {
                    "type": "string",
                    "description": "作者"
                },
                "repository": {
                    "type": "string",
                    "description": "仓库地址"
                }
            }
        },
        "UpdatePluginRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "description": "插件名称"
                },
                "version": {
                    "type": "string",
                    "description": "插件版本"
                },
                "description": {
                    "type": "string",
                    "description": "插件描述"
                },
                "author": {
                    "type": "string",
                    "description": "作者"
                },
                "repository": {
                    "type": "string",
                    "description": "仓库地址"
                }
            }
        },
        "PluginResponse": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string",
                    "description": "插件key"
                },
                "name": {
                    "type": "string",
                    "description": "插件名称"
                },
                "version": {
                    "type": "string",
                    "description": "插件版本"
                },
                "description": {
                    "type": "string",
                    "description": "插件描述"
                },
                "author": {
                    "type": "string",
                    "description": "作者"
                },
                "repository": {
                    "type": "string",
                    "description": "仓库地址"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time",
                    "description": "创建时间"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time",
                    "description": "更新时间"
                }
            }
        },
        "PluginListResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "description": "状态码"
                },
                "message": {
                    "type": "string",
                    "description": "状态信息"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "total": {
                            "type": "integer",
                            "description": "总数"
                        },
                        "items": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/PluginResponse"
                            }
                        }
                    }
                }
            }
        },
        "UploadResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "description": "状态码"
                },
                "message": {
                    "type": "string",
                    "description": "状态信息"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "filename": {
                            "type": "string",
                            "description": "文件名"
                        },
                        "size": {
                            "type": "integer",
                            "description": "文件大小"
                        },
                        "url": {
                            "type": "string",
                            "description": "文件URL"
                        }
                    }
                }
            }
        },
        "HealthResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "description": "状态码"
                },
                "message": {
                    "type": "string",
                    "description": "状态信息"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "status": {
                            "type": "string",
                            "description": "系统状态"
                        },
                        "version": {
                            "type": "string",
                            "description": "系统版本"
                        },
                        "timestamp": {
                            "type": "string",
                            "format": "date-time",
                            "description": "时间戳"
                        }
                    }
                }
            }
        },
        "Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "description": "状态码"
                },
                "message": {
                    "type": "string",
                    "description": "状态信息"
                }
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "description": "错误码"
                },
                "message": {
                    "type": "string",
                    "description": "错误信息"
                }
            }
        }
    }
}` 