package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"xobotyi.github.io/go/go-vanity-ssg/internal/config"
)

func TestPackage_VersionedPackages(t *testing.T) {
	t.Parallel()

	pkg := func(name, description string, versions []int) config.Package {
		return config.Package{
			Name:          name,
			Description:   description,
			Source:        nil,
			PrivateSource: nil,
			Versions:      versions,
		}
	}

	tests := []struct {
		name     string
		pkg      config.Package
		expected []config.Package
	}{
		{
			name:     "package without versions",
			pkg:      pkg("test-package", "Test package", nil),
			expected: nil,
		},
		{
			name: "package with versions",
			pkg:  pkg("test-package", "Test package", []int{2, 3}),
			expected: []config.Package{
				pkg("test-package/v2", "Test package", nil),
				pkg("test-package/v3", "Test package", nil),
			},
		},
		{
			name: "package with single version",
			pkg:  pkg("another-package", "Another test", []int{2}),
			expected: []config.Package{
				pkg("another-package/v2", "Another test", nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.pkg.VersionedPackages()
			assert.Equal(t, tt.expected, result)
		})
	}
}
