package cmd

import (
	"dev.gaijin.team/go/golib/must"
	"github.com/spf13/cobra"

	"github.com/xobotyi/go-vanity-ssg/internal/template"
)

func newEmitTemplatesCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emit-templates",
		Short: "Emit template files.",

		Args: cobra.NoArgs,

		RunE: emitTemplatesRunE,
	}

	fs := cmd.Flags()

	fs.String("dir", "./templates", "Directory to emit template files to.")
	must.NoErr(cmd.MarkFlagDirname("dir"))

	fs.Bool("overwrite", false, "Overwrite existing files.")

	return cmd
}

func emitTemplatesRunE(cmd *cobra.Command, _ []string) error {
	fs := cmd.Flags()

	dir := must.OK(fs.GetString("dir"))
	overwrite := must.OK(fs.GetBool("overwrite"))

	if err := ensureDir(dir); err != nil {
		return err
	}

	return template.WriteTemplatesDir(dir, overwrite, 0644)
}
