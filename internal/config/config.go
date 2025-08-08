package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"dev.gaijin.team/go/golib/e"
)

type Config struct {
	OutDir       string `yaml:"out-dir"`
	TemplatesDir string `yaml:"templates-dir,omitempty"`

	VanityRoot string `yaml:"vanity-root"`

	Packages Packages `yaml:"packages"`
}

type Package struct {
	Name          string         `yaml:"name"`
	Description   string         `yaml:"description"`
	Source        *PackageSource `yaml:"source"`
	PrivateSource *PackageSource `yaml:"private-source"`
	Versions      []int          `yaml:"versions,omitempty"`
}

type PackageSource struct {
	VcsType string   `yaml:"vcs-type"`
	VcsURI  string   `yaml:"vcs-uri"`
	URI     string   `yaml:"uri"`
	DirURI  string   `yaml:"dir-uri"`
	FileURI string   `yaml:"file-uri"`
	Swag    []string `yaml:"swag,omitempty"`
}

// VersionedPackages returns versioned package entries if versions are defined.
func (p Package) VersionedPackages() []Package {
	if len(p.Versions) == 0 {
		return nil
	}

	result := make([]Package, 0, len(p.Versions))
	for _, version := range p.Versions {
		versionedPkg := p
		versionedPkg.Name = fmt.Sprintf("%s/v%d", p.Name, version)
		versionedPkg.Versions = nil
		result = append(result, versionedPkg)
	}

	return result
}

type Packages []Package

// Public retrieves a list of packages that only have public source defined.
func (p Packages) Public() []Package {
	result := make([]Package, 0, len(p))

	for _, pkg := range p {
		if pkg.Source != nil {
			pkg.PrivateSource = nil
			result = append(result, pkg)
		}
	}

	return result
}

// Private returns a list of private packages with its Source fields replaced
// with private definition. In case withPublic is set to true - method also
// returns packages that contain public source.
func (p Packages) Private(withPublic bool) []Package {
	result := make([]Package, 0, len(p))

	for _, pkg := range p {
		if pkg.PrivateSource != nil {
			pkg.Source = pkg.PrivateSource
			result = append(result, pkg)

			continue
		}

		if withPublic && pkg.Source != nil {
			result = append(result, pkg)
		}
	}

	return result
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
