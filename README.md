# Gonfig

Gonfig is a Go library for managing application configuration files with support for JSON and YAML formats, file watching, and validation. It provides a generic interface for loading, saving, and validating configuration structs, and automatically reloads configuration when the file changes.

## Features

- Load and save configuration from JSON or YAML files
- Watch for file changes and reload configuration automatically
- Add custom validation functions for your config
- Thread-safe file operations

## Installation

```sh
go get github.com/michaelladouceur1/gonfig
```

## Usage

### 1. Define Your Config Struct

```go
type ConfigUserInfo struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

type Config struct {
    AppName  string         `json:"app_name"`
    Port     int            `json:"port"`
    UserInfo ConfigUserInfo `json:"user_info"`
}
```

### 2. Create a Validator (Optional)

```go
func validator(c *Config) error {
    if c.Port <= 0 || c.Port > 65535 {
        return &gonfig.ValidationError{Field: "Port", Message: "must be between 1 and 65535"}
    }
    return nil
}
```

### 3. Initialize Gonfig

```go
c, err := gonfig.NewGonfig(&Config{
    AppName: "my app",
    Port:    8080,
    UserInfo: ConfigUserInfo{
        Username: "admin",
        Email:    "admin@example.com",
    },
}, gonfig.GonfigFileOptions{
    Type:    gonfig.YAML, // or gonfig.JSON and gonfig.TOML
    RootDir: ".",
    Name:    "config",
    Watch:   true, // enable file watching
})

if err != nil {
    // handle error
}
```

### 4. Add Validators

```go
c.AddValidator(validator)
```

### 5. Update, Save, and Print Config

```go
c.Update(&Config{
    AppName: "updated app",
    Port:    8081,
    UserInfo: ConfigUserInfo{
        Username: "newadmin",
        Email:    "newadmin@example.com",
    },
})

if err := c.Save(); err != nil {
    // handle error
}

if err := c.PrintConfig(); err != nil {
    // handle error
}
```

### 6. Manual Validation

```go
c.Config.Port = 70000 // invalid port

if err := c.Validate(); err != nil {
    // handle error
}

if err := c.Save(); err != nil {
    // handle error
}
```

### 7. Watch for Changes

If `Watch` is enabled, Gonfig will reload the config and re-validate automatically when the file changes.

## API Reference

- `Gonfig`: Main generic config manager.
- `NewGonfig`: Create a new Gonfig instance.
- `GonfigFileOptions`: Options for config file type, location, and watching.
- `ValidationError`: Error type for validation failures.

## Supported File Types

- JSON (`gonfig.JSON`)
- YAML (`gonfig.YAML`)
- TOML (`gonfig.TOML`)

## License

MIT

---

See main.go for a complete example.  
See gonfig.go for implementation details.
