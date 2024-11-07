package main

import (
	"fmt"
	"os"

	"github.com/xobotyi/go-vanity-ssg/internal/cmd"
)

func main() {
	err := cmd.NewRootCMD().Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	os.Exit(0)
}
