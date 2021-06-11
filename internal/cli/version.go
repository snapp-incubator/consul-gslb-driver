package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(c *Config) *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(c *cobra.Command, args []string) {
			version()
		},
	}
	return versionCmd
}

func version() {
	fmt.Println("0.1.0")
}
