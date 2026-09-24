# oconf

A Go library for loading configuration files in CJSON format (JSON with C-style comments) with support for file inclusion and public/private annotations.

## Features

- **C-style comments**: Support for both single-line (`//`) and multi-line (`/* */`) comments in JSON files
- **File inclusion**: Use `#include` directive to merge multiple configuration files, at root level or inside any nested object
- **Public/private annotations**: Mark configuration values as public using `#public` for selective exposure
- **Multiple output formats**: Nested maps, flat maps (dot notation), and JSON strings
- **Thread-safe**: No global state - safe for concurrent use
- **Proper error handling**: All functions return errors instead of panicking
- **Type-safe**: Safe type assertions with proper error messages

## Installation

```bash
go get github.com/muqsith/oconf
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/muqsith/oconf"
)

func main() {
    // Load configuration
    config, err := oconf.GetConfig("config.cjson")
    if err != nil {
        log.Fatal(err)
    }

    // Access nested values
    server := config["server"].(map[string]interface{})
    port := server["port"].(float64)
    fmt.Printf("Server port: %.0f\n", port)
}
```

### Example Configuration File

Create a file named `config.cjson`:

```json
{
    // This is a comment
    "#include": ["base.cjson"],

    /*
     * Multi-line comment
     * Server configuration
     */
    "server": {
        "host#public": "example.com",  // Public values are marked with #public
        "port#public": 8080,
        "apiKey": "secret123"           // Private value - not marked public
    },

    "database": {
        "host": "localhost",
        "user": "admin",
        "password": "secret"
    }
}
```

## API Reference

### GetConfig

Load a configuration file and return it as a nested map.

```go
config, err := oconf.GetConfig("config.cjson")
if err != nil {
    log.Fatal(err)
}
```

### GetFlatConfig

Load a configuration file and return it as a flat map where nested keys are joined with dots.

```go
config, err := oconf.GetFlatConfig("config.cjson")
if err != nil {
    log.Fatal(err)
}

// Access values with dot notation
host := config["database.host"].(string)
port := config["server.port"].(float64)
```

### GetConfigAsJSONString

Load a configuration file and return it as a JSON string.

```go
jsonStr, err := oconf.GetConfigAsJSONString("config.cjson")
if err != nil {
    log.Fatal(err)
}
fmt.Println(jsonStr)
```

### GetPublicConfig

Load only values marked with `#public` annotation. Useful for exposing configuration to external systems while keeping sensitive values private.

```go
publicConfig, err := oconf.GetPublicConfig("config.cjson")
if err != nil {
    log.Fatal(err)
}

// Only returns: {"server": {"host": "example.com", "port": 8080}}
// Database credentials are excluded
```

### GetFlatPublicConfig

Get public configuration as a flat map.

```go
publicConfig, err := oconf.GetFlatPublicConfig("config.cjson")
if err != nil {
    log.Fatal(err)
}

// Access with dot notation
host := publicConfig["server.host"].(string)
```

### GetPublicConfigAsJSONString

Get public configuration as a JSON string.

```go
jsonStr, err := oconf.GetPublicConfigAsJSONString("config.cjson")
if err != nil {
    log.Fatal(err)
}

// Send to external API or client
```

## Advanced Features

### File Inclusion

Use the `#include` directive to merge multiple configuration files. Included files are merged in order, and values in the current file override included values.

**base.cjson:**
```json
{
    "server": {
        "host": "localhost",
        "port": 8000,
        "timeout": 30
    }
}
```

**development.cjson:**
```json
{
    "#include": ["base.cjson"],
    "server": {
        "port": 8080  // Overrides port from base.cjson
    }
}
```

Result: `server.host` is "localhost" (from base), `server.port` is 8080 (overridden), and `server.timeout` is 30 (from base).

### Nested Includes

`#include` works at any depth inside a config file, not just at the root. This is useful for composing specific sections from separate files.

**prod.cjson:**
```json
{
    "host": "db.prod.example.com",
    "port": 5432
}
```

**secrets.cjson:**
```json
{
    "password": "s3cr3t",
    "user": "app_user"
}
```

**config.cjson:**
```json
{
    "app": "myapp",
    "database": {
        "#include": [
            "./profile-service/prod.cjson",
            "./profile-service/secrets.cjson"
        ],
        "name": "mydb"
    }
}
```

Result:
```json
{
    "app": "myapp",
    "database": {
        "host": "db.prod.example.com",
        "port": 5432,
        "password": "s3cr3t",
        "user": "app_user",
        "name": "mydb"
    }
}
```

**Rules for nested includes:**
- Multiple files are merged left-to-right; later files win on key conflicts.
- Local keys in the nested object override keys from any included file.
- A nested `#include` belongs to the current file, so its values override values for the same key pulled in by a root-level `#include`.
- The same file can be included from several nested objects; each one gets the full contents. Only circular includes (a file including itself through a chain) are skipped.
- Include paths are relative to the directory of the file that contains the `#include`.
- Included files can themselves use `#include` at their own root level.
- Errors (missing file, wrong type) surface immediately with the failing key name in the message.

### Relative and Absolute Paths

Include paths can be:
- Relative to the current config file: `"#include": ["./configs/base.cjson"]`
- Absolute paths: `"#include": ["/etc/app/base.cjson"]`

### Public/Private Annotations

Mark keys with `#public` to expose them selectively:

```json
{
    "server#public": {
        "host#public": "example.com",
        "apiKey": "secret"  // Not public
    },
    "database": {
        "password": "secret123"  // Not public
    }
}
```

When using `GetPublicConfig()`, only values marked with `#public` are returned, and the `#public` suffix is removed from keys.

## Error Handling

All functions return proper errors. Common error types:

- **File errors**: File not found, permission denied, etc.
- **JSON errors**: Invalid JSON syntax
- **Type errors**: Unexpected value types in configuration

```go
config, err := oconf.GetConfig("config.cjson")
if err != nil {
    switch e := err.(type) {
    case *oconf.FileError:
        fmt.Printf("File error: %v\n", e)
    case *oconf.TypeError:
        fmt.Printf("Type error: %v\n", e)
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Type Assertions

When accessing configuration values, use type assertions with the comma-ok idiom for safety:

```go
config, err := oconf.GetConfig("config.cjson")
if err != nil {
    log.Fatal(err)
}

// Safe type assertion
if server, ok := config["server"].(map[string]interface{}); ok {
    if port, ok := server["port"].(float64); ok {
        fmt.Printf("Port: %.0f\n", port)
    }
}
```

## Testing

Run the tests:

```bash
# Run all tests
go test

# Run specific test
go test -run TestGetConfig1

# Run with verbose output
go test -v
```

## Version History

### v1.1.0 (Latest)
- ✅ `#include` now supported inside nested objects at any depth

### v1.0.0
- ✅ Proper error handling - all functions return errors
- ✅ Thread-safe - no global state
- ✅ Safe type assertions
- ✅ Comprehensive documentation
- ✅ Bug fixes and performance improvements
- ✅ Improved path handling for cross-platform compatibility

### v0.x (Legacy)
- Basic functionality with `log.Fatal` error handling
- Global state for file tracking

## Migration from v0.x to v1.0

If you're upgrading from an older version, you need to handle errors:

**Old code:**
```go
config := oconf.GetConfig("config.cjson")
```

**New code:**
```go
config, err := oconf.GetConfig("config.cjson")
if err != nil {
    log.Fatal(err)
}
```

## Publishing a New Version

### 1. Run the tests

```bash
go test ./... -v -count=1
```

### 2. Tag and push

```bash
git fetch --all --tags   # sync existing tags
git tag -l               # verify the next version
git tag v1.x.y
git push origin v1.x.y
```

### 3. Force the proxy to index the new version

```bash
GOPROXY=proxy.golang.org go list -m github.com/muqsith/oconf@v1.x.y
```

Then verify the release is live at https://pkg.go.dev/github.com/muqsith/oconf.

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is a near port of the [oconf npm package](https://www.npmjs.com/package/oconf).

## Things Not Yet Supported

- Relaxed JSON (JSON5)
- Error on conflict when multiple files define the same `#public` attribute

## Acknowledgments

This is a Go port of the [oconf](https://www.npmjs.com/package/oconf) JavaScript library.
