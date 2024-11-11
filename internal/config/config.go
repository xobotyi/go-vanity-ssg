package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"dev.gaijin.team/go/golib/e"
)

type Config struct {
	OutDir       string    `yaml:"out-dir"`
	TemplatesDir string    `yaml:"templates-dir"`
	Packages     []Package `yaml:"packages"`
}

type Package struct {
	Name          string        `yaml:"name"`
	Description   string        `yaml:"description"`
	Source        PackageSource `yaml:"source"`
	PrivateSource PackageSource `yaml:"private-source"`
}

type PackageSource struct {
	VcsType string `yaml:"vcs-type"`
	VcsURI  string `yaml:"vcs-uri"`
	URI     string `yaml:"uri"`
	DirURI  string `yaml:"dir-uri"`
	FileURI string `yaml:"file-uri"`
}

// Parse config file.
func Parse(filePath string) (Config, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, e.NewFrom("failed to read config file", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, e.NewFrom("failed to parse config file", err)
	}

	return cfg, nil
}
