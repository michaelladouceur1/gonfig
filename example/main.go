package main

import "github.com/michaelladouceur1/gonfig"

type AppConfig struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Server      struct {
		Host string `json:"host" yaml:"host"`
		Port int    `json:"port" yaml:"port"`
	} `json:"server" yaml:"server"`
}

func validator(config AppConfig) error {
	if config.Name == "" {
		return &gonfig.ValidationError{Field: "Name", Message: "Name cannot be empty"}
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return &gonfig.ValidationError{Field: "Server.Port", Message: "Port must be between 1 and 65535"}
	}
	return nil
}

func main() {
	cfg := AppConfig{
		Name:        "MyApp",
		Description: "This is my application",
		Server: struct {
			Host string `json:"host" yaml:"host"`
			Port int    `json:"port" yaml:"port"`
		}{
			Host: "localhost",
			Port: 8080,
		},
	}

	opts := gonfig.GonfigFileOptions{
		Type:    gonfig.YAML,
		RootDir: ".",
		Name:    "config",
		Watch:   true,
	}

	config, err := gonfig.NewGonfig(cfg, opts)
	if err != nil {
		panic(err)
	}

	config.AddValidator(validator)

	config.Config.Server.Port = 9090

	if err := config.Validate(); err != nil {
		panic(err)
	}

	if err := config.Save(); err != nil {
		panic(err)
	}

	config.PrintConfig()
}
