package cmd

import (
	"github.com/spf13/cobra"

	"dev.gaijin.team/go/golib/must"

	"github.com/xobotyi/go-vanity-ssg/internal/config"
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

			cfg, err := config.Parse(must.OK(fs.GetString("config")))
			if err != nil {
				return err //nolint:wrapcheck
			}

			if cfg.OutDir == "" || outDir != DefaultOutDir {
				cfg.OutDir = outDir
			}

			return nil
		},
	}

	fs := cmd.Flags()

	fs.StringP("out-dir", "o", DefaultOutDir, "Directory to emit html files.")
	must.NoErr(cmd.MarkFlagDirname("out-dir"))

	fs.StringP("config", "c", DefaultConfigPath, "Path to config file (.yaml, .yml).")
	must.NoErr(cmd.MarkFlagFilename("config", ".yaml", ".yml"))

	return cmd
}
