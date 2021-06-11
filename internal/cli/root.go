package cli

import (
	goflag "flag"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	cfgFile string
	config  = &Config{}

	rootCmd = &cobra.Command{
		Use:   "consul-gslb-driver",
		Short: "Gslb driver for Hashicorp Consul",
		Long:  `Gslb driver for Hashicorp Consul to run as a side-car of gslb-controller`,
	}
)

func init() {

	rootCmd.Flags().SortFlags = false
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		newVersionCmd(config),
		newStartCmd(config),
	)
	klog.InitFlags(nil)
	goflag.Parse()
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is config.yaml)")

	//// local command
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.ErrorS(err, "rootCmd.Execute()")
		os.Exit(1)
	}
}

func initConfig() {
	// map command line flags to viper variables.
	// viper will prefer flags from command line rather than file
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		klog.ErrorS(err, "Failed bind flags")
	}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		klog.InfoS("Using config", "file", viper.ConfigFileUsed())
	} else {
		klog.ErrorS(err, "Failed to read config")
	}
	if err := viper.Unmarshal(config); err != nil {
		klog.ErrorS(err, "Failed to unmarshal config")
	}
	verbosity := viper.GetString("v")
	// Manually update v flag from viper config file
	if verbosity == "0" {
		if err := goflag.Set("v", viper.GetString("v")); err != nil {
			klog.Errorf("%+v", err)
			return
		}
	}
	klog.InfoS("Verbosity", "v", verbosity)
}
