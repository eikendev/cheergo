// Package commands contains functions that are exposed as dedicated commands of the tool.
package commands

import (
	"fmt"
	"runtime/debug"
)

// VersionCommand represents the options specific to the version command.
type VersionCommand struct{}

// Run is the function for the version command.
func (*VersionCommand) Run() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("build info not available")
	}

	fmt.Printf("cheergo %s\n", buildInfo.Main.Version)
	return nil
}
