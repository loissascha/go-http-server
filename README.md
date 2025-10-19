# go-http-server

A lightweight and flexible HTTP server framework for Go, designed to simplify building web applications with built-in support for routing, middleware, internationalization, CORS, and OpenAPI documentation generation.

## Features

- **Simple Routing**: Support for GET, POST, PUT, and DELETE methods with clean route definitions
- **Middleware Support**: Chain multiple middlewares for request processing
- **Internationalization (i18n)**: Built-in translation support with JSON files, auto-language detection from URLs, and fallback to default language
- **CORS Handling**: Automatic CORS middleware for cross-origin requests
- **Panic Recovery**: Middleware to recover from panics and prevent server crashes
- **OpenAPI Generation**: Automatically generate OpenAPI 3.1.0 JSON specifications from your routes
- **Response Helpers**: Convenient JSON response functions
- **Server Options**: Configurable server setup with various options

## Installation

```bash
go get github.com/loissascha/go-http-server
```

## Basic Usage

```go
package main

import (
    "net/http"
    "github.com/loissascha/go-http-server/respond"
    "github.com/loissascha/go-http-server/server"
)

func main() {
    // Create server with options
    s, err := server.NewServer(
        server.EnableTranslations(),
        server.EnableAutoDetectLanguage(),
        server.AddTranslationFile("en", "en.json"),
        server.AddTranslationFile("de", "de.json"),
        server.SetDefaultLanguage("en"),
    )
    if err != nil {
        panic(err)
    }

    // Define routes
    s.GET("/", homeHandler)
    s.POST("/api/data", dataHandler)

    // Start server
    s.Serve(":8080")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    // Get translated string
    message := s.GetLanguageString(r, "welcome")
    respond.JSON(w, http.StatusOK, map[string]string{"message": message})
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    respond.JSON(w, http.StatusOK, map[string]string{"status": "success"})
}
```

## Translation Files

Create JSON files for translations (e.g., `en.json`):

```json
{
  "welcome": "Welcome to our API",
  "error": "An error occurred"
}
```

## OpenAPI Documentation

Generate OpenAPI JSON specification:

```go
s.CreateOpenAPIJson("8080") // Creates openapi.json
```
