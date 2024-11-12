package template

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"net/url"
	"os"
	"path"
	"strings"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/must"

	"github.com/xobotyi/go-vanity-ssg/internal/config"
)

type Vanity struct {
	vanityRoot string
	tpl        *template.Template
}

func New(vanityRoot string) Vanity {
	return Vanity{
		vanityRoot: vanityRoot,
		tpl: template.New("_").Funcs(map[string]any{
			"unescapeHTML": func(s string) template.HTML {
				return template.HTML(s) //nolint:gosec
			},
		}),
	}
}

//go:embed templates/*.gohtml
var tpls embed.FS

const (
	indexTemplate   = "index.gohtml"
	packageTemplate = "package.gohtml"
)

var templatesNames = []string{indexTemplate, packageTemplate} //nolint:gochecknoglobals

func buildTemplatePaths(dir string) []string {
	globs := make([]string, 0, len(templatesNames))

	for _, name := range templatesNames {
		globs = append(globs, path.Join(dir, name))
	}

	return globs
}

func fileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return !os.IsNotExist(err)
}

// WriteTemplatesDir writes all embedded templates to output directory.
func WriteTemplatesDir(dir string, overwrite bool, perms os.FileMode) error {
	for _, tname := range templatesNames {
		tpath := path.Join(dir, tname)
		if fileExists(tpath) && !overwrite {
			return e.New("target template file already exists", fields.F("path", tpath))
		}

		tmplContent := must.OK(fs.ReadFile(tpls, path.Join("templates", tname)))

		if err := os.WriteFile(tpath, tmplContent, perms); err != nil {
			return e.NewFrom("failed to write file", err, fields.F("path", tpath))
		}
	}

	return nil
}

// ParseTemplates from embedded fs and provided override directory.
//
// Override directory may contain subset of templates. It only expects `.gohtml`
// files.
func (v *Vanity) ParseTemplates(overrideDir string) (err error) {
	v.tpl, err = v.tpl.ParseFS(tpls, buildTemplatePaths("templates")...)
	if err != nil {
		return e.NewFrom("failed to parse template", err)
	}

	if overrideDir != "" {
		for _, g := range buildTemplatePaths(overrideDir) {
			tt, err := v.tpl.ParseGlob(g)

			// ParseGlob returns error that cant be checked with errors.Is
			// therefore we have to check by substring
			if err != nil && !strings.Contains(err.Error(), "pattern matches no files") {
				return err
			}

			if tt != nil {
				v.tpl = tt
			}
		}
	}

	v.tpl = v.tpl.Funcs(map[string]any{
		"unescapeHTML": func(s string) template.HTML {
			return template.HTML(s) //nolint:gosec
		},
	})

	return nil
}

type indexData struct {
	Title    string
	Packages []packageData
}

func (v *Vanity) EmitIndex(outDir string, pkgs []config.Package) error {
	b, err := v.renderIndex(pkgs)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(outDir, "index.html"), b, 0644) //nolint:mnd
	if err != nil {
		return e.NewFrom("failed to emit index file", err)
	}

	return nil
}

func (v *Vanity) renderIndex(pkgs []config.Package) ([]byte, error) {
	buf := &bytes.Buffer{}

	id := indexData{
		Title:    v.vanityRoot,
		Packages: make([]packageData, 0, len(pkgs)),
	}

	for _, p := range pkgs {
		pd, err := v.packageData(p)
		if err != nil {
			return nil, err
		}

		id.Packages = append(id.Packages, pd)
	}

	if err := v.tpl.ExecuteTemplate(buf, indexTemplate, id); err != nil {
		return nil, e.NewFrom("failed to render index template", err)
	}

	return buf.Bytes(), nil
}

// EmitPackages render all packages from provided config and writes them into dir
// folder.
func (v *Vanity) EmitPackages(outDir string, pkgs []config.Package) error {
	for _, pkg := range pkgs {
		b, err := v.renderPackage(pkg)
		if err != nil {
			return err
		}

		err = os.WriteFile(path.Join(outDir, path.Base(pkg.Name)+".html"), b, 0644) //nolint:mnd
		if err != nil {
			return e.NewFrom("failed to emit package file", err)
		}
	}

	return nil
}

type packageData struct {
	FQN string
	config.Package
}

func (v *Vanity) packageData(p config.Package) (packageData, error) {
	packageFQN, err := url.JoinPath(v.vanityRoot, p.Name)
	if err != nil {
		return packageData{}, e.NewFrom("failed to build package fully qualified name", err)
	}

	return packageData{
		FQN:     packageFQN,
		Package: p,
	}, nil
}

func (v *Vanity) renderPackage(p config.Package) ([]byte, error) {
	buf := &bytes.Buffer{}

	pd, err := v.packageData(p)
	if err != nil {
		return nil, err
	}

	if err := v.tpl.ExecuteTemplate(buf, packageTemplate, pd); err != nil {
		return nil, e.NewFrom("failed to render package template", err)
	}

	return buf.Bytes(), nil
}
