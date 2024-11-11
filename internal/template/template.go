package template

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path"
	"strings"

	"dev.gaijin.team/go/golib/e"

	"github.com/xobotyi/go-vanity-ssg/internal/config"
)

type Vanity struct {
	pkgs []config.Package
	tpl  *template.Template
}

func New(pkgs []config.Package) Vanity {
	return Vanity{pkgs: pkgs}
}

//go:embed templates/*.gohtml
var tpls embed.FS

const (
	indexTemplate   = "index.gohtml"
	packageTemplate = "package.gohtml"
)

var templatesNames = []string{indexTemplate, packageTemplate}

func buildTemplatePaths(dir string) []string {
	globs := make([]string, 0, len(templatesNames))

	for _, name := range templatesNames {
		globs = append(globs, path.Join(dir, name))
	}

	return globs
}

// ParseTemplates from embedded fs and provided override directory.
//
// Override directory may contain subset of templates. It only expects `.gohtml`
// files.
func (v *Vanity) ParseTemplates(overrideDir string) (err error) {
	v.tpl, err = template.ParseFS(tpls, buildTemplatePaths("templates")...)
	if err != nil {
		return err
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

	return nil
}

type indexData struct {
	Title    string
	Packages []indexPackageData
}

type indexPackageData struct {
	URI  string
	Name string
}

func (v *Vanity) EmitIndex(outDir string) error {
	id := indexData{
		Title:    "TEST TITLE",
		Packages: nil,
	}

	for _, pkg := range v.pkgs {
		id.Packages = append(id.Packages, indexPackageData{
			URI:  "",
			Name: pkg.Name,
		})
	}

	buf := &bytes.Buffer{}

	if err := v.tpl.ExecuteTemplate(buf, indexTemplate, id); err != nil {
		return e.NewFrom("failed to render index template", err)
	}

	err := os.WriteFile(path.Join(outDir, "index.html"), buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vanity) renderIndex() ([]byte, error) {
	return nil, nil
}

// EmitPackages render all packages from provided config and writes them into dir
// folder.
func (v *Vanity) EmitPackages(outDir string) error {
	for _, pkg := range v.pkgs {
		b, err := v.renderPackage(pkg)
		if err != nil {
			return err
		}

		err = os.WriteFile(path.Join(outDir, path.Base(pkg.Name)+".html"), b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Vanity) renderPackage(p config.Package) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := v.tpl.ExecuteTemplate(buf, packageTemplate, p); err != nil {
		return nil, e.NewFrom("failed to render package template", err)
	}

	return buf.Bytes(), nil
}
