package cli

import (
	"github.com/snapp-incubator/consul-gslb-driver/internal/consul"
	"github.com/snapp-incubator/consul-gslb-driver/internal/servers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func newStartCmd(c *Config) *cobra.Command {

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the driver",
		Run: func(cmd *cobra.Command, args []string) {
			start(c)
		},
	}
	var flag string

	flag = "grpcAddress"
	startCmd.Flags().StringVar(&c.GRPCAddress, flag, "unix://var/run/gslbi/gslbi.sock", "grpc address to listen on")
	if err := viper.BindPFlag(flag, startCmd.Flags().Lookup(flag)); err != nil {
		klog.ErrorS(err, "unable to bind flag")
	}

	flag = "metrics-port"
	startCmd.Flags().IntVarP(&c.MetricServer.Port, flag, "p", 8080, "port to bind to")
	if err := viper.BindPFlag(flag, startCmd.Flags().Lookup(flag)); err != nil {
		klog.ErrorS(err, "unable to bind flag")
	}

	return startCmd
}

func start(c *Config) {
	d := servers.NewDriver(c.GRPCAddress, c.MetricServer.IP, c.MetricServer.Path, c.MetricServer.Port)
	// Initialize consul with passed configuration
	consul, err := consul.GetConsul(c.ConsulConfig.Address, c.ConsulConfig.Scheme, c.ConsulConfig.Datacenter)
	if err != nil {
		klog.Warningf("Failed to GetConsul: %v", err)
		return
	}
	d.SetupDriver(consul)
	d.Run()

	// addr := net.JoinHostPort(c.Listen.IP, strconv.Itoa(c.Listen.Port))
	// _, cancel := context.WithCancel(context.Background())
	// servers.RunServer(cancel, addr)

}
