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

## Environment Variables

This server currently uses environment variables for runtime environment behavior and CORS origin validation.

| Variable | Type | Default | Description |
| --- | --- | --- | --- |
| `APP_ENV` | enum-like string (`production` or non-production) | empty (treated as non-production) | Controls strictness for CORS behavior. Only the exact value `production` enables strict production checks. |
| `ALLOWED_ORIGINS` | comma-separated string, or `*` | if `APP_ENV != production`: built-in localhost allowlist; if `APP_ENV == production`: required | List of allowed browser origins for CORS (for example `https://app.example.com,https://admin.example.com`). |

### `APP_ENV`

- Use `APP_ENV=production` in production deployments.
- Any other value (or not setting it) is treated as non-production.
- In non-production, localhost origins are allowed automatically as a developer convenience.

### `ALLOWED_ORIGINS`

- Accepts a comma-separated list of origins (spaces are allowed and trimmed).
- Example: `ALLOWED_ORIGINS=https://app.example.com, https://admin.example.com`
- You can set `ALLOWED_ORIGINS=*` to allow any origin (use with care).
- In production (`APP_ENV=production`), this should always be set to explicit trusted origins.
- If `APP_ENV=production` and `ALLOWED_ORIGINS` is missing or invalid, the server fails fast during startup.

### Recommended `.env` examples

Development:

```env
APP_ENV=development
# Optional in development; localhost origins are auto-allowed when unset
# ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

Production:

```env
APP_ENV=production
ALLOWED_ORIGINS=https://app.example.com,https://admin.example.com
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
