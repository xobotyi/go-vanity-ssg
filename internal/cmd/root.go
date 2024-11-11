package cmd

import (
	"errors"
	"os"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"github.com/spf13/cobra"

	"dev.gaijin.team/go/golib/must"

	"github.com/xobotyi/go-vanity-ssg/internal/config"
	"github.com/xobotyi/go-vanity-ssg/internal/template"
)

const (
	DefaultOutDir     = "./dist"
	DefaultConfigPath = "./.vanity.config.yaml"
)

func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "go-vanity",
		Short: "Vanity imports static site generator.",

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, _ []string) error {
			fs := cmd.Flags()

			outDir := must.OK(fs.GetString("out-dir"))
			templatesDir := must.OK(fs.GetString("templates-dir"))

			cfg, err := config.Parse(must.OK(fs.GetString("config")))
			if err != nil {
				return err //nolint:wrapcheck
			}

			if cfg.OutDir == "" || outDir != DefaultOutDir {
				cfg.OutDir = outDir
			}

			if err := ensureDir(cfg.OutDir); err != nil {
				return err
			}

			if templatesDir != "" {
				cfg.TemplatesDir = templatesDir
			}

			if cfg.TemplatesDir != "" {
				if err := pathIsDir(cfg.TemplatesDir); err != nil {
					return err
				}
			}

			vt := template.New(cfg.Packages)

			err = vt.ParseTemplates(cfg.TemplatesDir)
			if err != nil {
				return err
			}

			if err := vt.EmitPackages(cfg.OutDir); err != nil {
				return err
			}

			if err := vt.EmitIndex(cfg.OutDir); err != nil {
				return err
			}

			return nil
		},
	}

	fs := cmd.Flags()

	fs.StringP("out-dir", "o", DefaultOutDir, "Directory to emit html files.")
	must.NoErr(cmd.MarkFlagDirname("out-dir"))

	fs.StringP("templates-dir", "t", "", "Directory containing override templates.")
	must.NoErr(cmd.MarkFlagDirname("templates-dir"))

	fs.StringP("config", "c", DefaultConfigPath, "Path to config file (.yaml, .yml).")
	must.NoErr(cmd.MarkFlagFilename("config", ".yaml", ".yml"))

	return cmd
}

func pathIsDir(path string) error {
	finfo, err := os.Stat(path)
	if err != nil {
		return e.NewFrom("unable to stat path", err)
	}

	if !finfo.IsDir() {
		return e.New("provided path is not a directory", fields.F("path", path))
	}

	return nil
}

func ensureDir(path string) error {
	finfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(path, 0644); err != nil {
				return e.NewFrom("failed to create output dir", err)
			}

			return nil
		}

		return e.NewFrom("unable to stat path", err)
	}

	if !finfo.IsDir() {
		return e.New("out path is not a directory")
	}

	return nil
}
