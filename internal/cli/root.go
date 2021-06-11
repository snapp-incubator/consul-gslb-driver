package cli

import (
	"github.com/spf13/cobra"
)

func newRootCommand(cfg *Config) *cobra.Command {

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "consul-gslb-driver",
		Short: "Gslb driver for Hashicorp Consul",
		Long:  `Gslb driver for Hashicorp Consul to run as a side-car of gslb-controller`,
		// PersistentPreRunE: func(c *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		// return initializeConfig(cfg)
		// },
	}
	rootCmd.AddCommand(newVersionCmd(cfg))
	rootCmd.AddCommand(newStartCmd(cfg))

	return rootCmd
}
