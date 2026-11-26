package template_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"xobotyi.github.io/go/go-vanity-ssg/internal/config"
	"xobotyi.github.io/go/go-vanity-ssg/internal/template"
)

func Test_Vanity(t *testing.T) {
	t.Parallel()

	t.Run("EmitPackages", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			giveSubdir string
			wantImport string
		}{
			{
				name:       "module at the repository root",
				giveSubdir: "",
				wantImport: `<meta name="go-import" content="go.example.com/pkg git https://example.com/repo">`,
			},
			{
				name:       "module in a subdirectory",
				giveSubdir: "bindings/go",
				wantImport: `<meta name="go-import" content="go.example.com/pkg git https://example.com/repo bindings/go">`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				outDir := t.TempDir()

				vanity := template.New("go.example.com")
				require.NoError(t, vanity.ParseTemplates(""))

				err := vanity.EmitPackages(outDir, []config.Package{{
					Name:        "pkg",
					Description: "",
					Source: &config.PackageSource{
						VcsType: "git",
						VcsURI:  "https://example.com/repo",
						Subdir:  tt.giveSubdir,
						URI:     "",
						DirURI:  "",
						FileURI: "",
						Swag:    nil,
					},
					PrivateSource: nil,
					Versions:      nil,
				}})
				require.NoError(t, err)

				page, err := fs.ReadFile(os.DirFS(outDir), "pkg.html")
				require.NoError(t, err)
				require.Contains(t, string(page), tt.wantImport)
			})
		}
	})
}
