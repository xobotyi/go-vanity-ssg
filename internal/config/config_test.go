package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"xobotyi.github.io/go/go-vanity-ssg/internal/config"
)

func TestPackage_VersionedPackages(t *testing.T) {
	tests := []struct {
		name     string
		pkg      config.Package
		expected []config.Package
	}{
		{
			name:     "package without versions",
			pkg:      config.Package{Name: "test-package", Description: "Test package"},
			expected: nil,
		},
		{
			name: "package with versions",
			pkg:  config.Package{Name: "test-package", Description: "Test package", Versions: []int{2, 3}},
			expected: []config.Package{
				{Name: "test-package/v2", Description: "Test package"},
				{Name: "test-package/v3", Description: "Test package"},
			},
		},
		{
			name: "package with single version",
			pkg:  config.Package{Name: "another-package", Description: "Another test", Versions: []int{2}},
			expected: []config.Package{
				{Name: "another-package/v2", Description: "Another test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkg.VersionedPackages()
			assert.Equal(t, tt.expected, result)
		})
	}
}
