package cmd

import (
	"os"
	"path"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/must"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"xobotyi.github.io/go/go-vanity-ssg/internal/config"
)

var cfgExample = config.Config{ //nolint:gochecknoglobals
	OutDir:       "./dist",
	TemplatesDir: "",
	VanityRoot:   "go.example.com",
	Packages: []config.Package{
		{
			Name:        "new-shiny-package",
			Description: "Description of a package",
			Source: &config.PackageSource{
				VcsType: "git",
				VcsURI:  "https://github.com/my/new-shiny-package",
				URI:     "https://github.com/my/new-shiny-package",
				DirURI:  "https://github.com/my/new-shiny-package{/dir}",
				FileURI: "https://github.com/my/new-shiny-package{/dir}/{file}#L{line}",
				Swag:    nil,
			},
			PrivateSource: &config.PackageSource{
				VcsType: "git",
				VcsURI:  "https://github.com/my/new-shiny-package",
				URI:     "https://github.com/my/new-shiny-package",
				DirURI:  "https://github.com/my/new-shiny-package{/dir}",
				FileURI: "https://github.com/my/new-shiny-package{/dir}/{file}#L{line}",
				Swag:    nil,
			},
		},
	},
}

func newEmitConfigCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emit-config",
		Short: "Emit example config file.",

		Args: cobra.NoArgs,

		RunE: emitConfigRunE,
	}

	fs := cmd.Flags()

	fs.String("config", DefaultConfigPath, "Path to config file.")
	must.NoErr(cmd.MarkFlagFilename("config", ".yaml", ".yml"))

	fs.Bool("overwrite", false, "Overwrite file if it exists.")
	must.NoErr(cmd.MarkFlagFilename("config", ".yaml", ".yml"))

	return cmd
}

func emitConfigRunE(cmd *cobra.Command, _ []string) error {
	fs := cmd.Flags()

	cfgPath := must.OK(fs.GetString("config"))
	overwrite := must.OK(fs.GetBool("overwrite"))

	if err := ensureDir(path.Dir(cfgPath)); err != nil {
		return err
	}

	if fileExists(cfgPath) && !overwrite {
		return e.New("target file already exists")
	}

	err := os.WriteFile(cfgPath, must.OK(yaml.Marshal(cfgExample)), 0644)
	if err != nil {
		return e.NewFrom("failed to write config example", err, fields.F("path", cfgPath))
	}

	return nil
}

func fileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return !os.IsNotExist(err)
}
