package cli

import (
	"context"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/snapp-cab/consul-gslb-driver/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStartCmd(c *Config) *cobra.Command {

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the driver",
		Run: func(cmd *cobra.Command, args []string) {
			start(c)
		},
	}
	startCmd.Flags().IntVarP(&c.Listen.Port, "port", "p", 8080, "port to bind to")
	if err := viper.BindPFlag("port", startCmd.Flags().Lookup("port")); err != nil {
		log.Fatal("Unable to bind flag:", err)
	}
	startCmd.Flags().IntVarP(&c.Verbose, "verbose", "v", 2, "verbosity")
	if err := viper.BindPFlag("verbose", startCmd.Flags().Lookup("verbose")); err != nil {
		log.Fatal("Unable to bind flag:", err)
	}
	return startCmd
}

func start(c *Config) {
	log.SetLevel(log.Level(c.Verbose))
	log.Info("Log level: ", c.Verbose)
	addr := net.JoinHostPort(c.Listen.IP, strconv.Itoa(c.Listen.Port))
	_, cancel := context.WithCancel(context.Background())
	server.RunServer(cancel, addr)

}
