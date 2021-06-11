package cli

import (
	"github.com/snapp-cab/consul-gslb-driver/internal/servers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	endpoint     string
	consulConfig string
	datacenter   string
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
		klog.ErrorS(err, "unable to bind flag")
	}
	return startCmd
}

func start(c *Config) {
	endpoint = c.Listen.IP
	d := servers.NewDriver(endpoint, datacenter)
	d.SetupDriver()
	d.Run()

	// addr := net.JoinHostPort(c.Listen.IP, strconv.Itoa(c.Listen.Port))
	// _, cancel := context.WithCancel(context.Background())
	// servers.RunServer(cancel, addr)

}
