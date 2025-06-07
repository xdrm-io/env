# env

Go package for reading environment variables with built-in support for Docker secrets.

## Key Feature

For any environment variable `VAR_NAME`, this package automatically checks:
1. `VAR_NAME` environment variable
2. `VAR_NAME_FILE` environment variable containing a file path

This makes it **Docker secrets compatible** out of the box.

## Installation

```bash
go get github.com/xdrm-io/env
```

## Basic Usage

### Simple Reading

```go
package main

import (
    "fmt"
    "github.com/xdrm-io/env"
)

func main() {
    // Reads DATABASE_URL or DATABASE_URL_FILE
    value, exists := env.Read("DATABASE_URL")
    if exists {
        fmt.Println("Database URL:", value)
    }
}
```

### Struct Decoding

```go
type Config struct {
    Host     string        `env:"DB_HOST,required"`
    Port     int           `env:"DB_PORT"`
    Password string        `env:"DB_PASSWORD,required"`
    Debug    bool          `env:"DEBUG"`
    Timeout  time.Duration `env:"TIMEOUT"`
    Tags     []string      `env:"TAGS"`
}

func main() {
    var config Config
    if err := env.ReadStruct(&config); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", config)
}
```

Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_PASSWORD_FILE=/run/secrets/db_password  # Docker secret
export DEBUG=true
export TIMEOUT=30s
export TAGS=web,api,production
```

## Docker Example

### docker-compose.yml
```yaml
services:
  app:
    image: myapp
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD_FILE=/run/secrets/db_password
    secrets:
      - db_password

secrets:
  db_password:
    file: ./db_password.txt
```

### Go code
```go
type AppConfig struct {
    DBHost     string `env:"DB_HOST,required"`
    DBPassword string `env:"DB_PASSWORD,required"`  // Reads from file
}

var config AppConfig
if err := env.ReadStruct(&config); err != nil {
    log.Fatal(err)
}
// config.DBPassword contains the content of /run/secrets/db_password
```

## Supported Types

- Basic: `string`, `[]byte`, `bool`
- Numbers: `int`, `int8/16/32/64`, `uint`, `uint8/16/32/64`, `float32/64`
- Time: `time.Time` (RFC3339), `time.Duration`
- Collections: `[]string` (comma-separated)
- Logging: `slog.Level` ("debug", "info", "warn", "error")

## Struct Tags

- `env:"VAR_NAME"` - binds field to environment variable
- `env:"VAR_NAME,required"` - makes field mandatory

## Error Types

```go
const (
    ErrNotPtr           // not a pointer
    ErrNotStructPtr     // not a pointer to struct
    ErrFieldRequired    // required field missing
    ErrFieldDecode      // decode error
    ErrFieldUnsupported // unsupported type
)
```

## Why Use This?

✅ **Secure**: Secrets never exposed in environment variables
✅ **Simple**: Same code for direct vars or file-based secrets
✅ **Compatible**: Works with Docker, Kubernetes, etc.
✅ **Type-safe**: Automatic decoding to Go types
✅ **Validated**: Required fields with clear error handling
