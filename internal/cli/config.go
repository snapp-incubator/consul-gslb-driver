package cli

// corespoding struct for config.yaml file
type Config struct {
	V            int          `mapstructure:"v"`
	GRPCAddress  string       `mapstructure:"grpcAddress"`
	ConsulConfig ConsulConfig `mapstructure:"consulConfig"`
	MetricServer MetricServer `mapstructure:"metricServer"`
}

type ConsulConfig struct {
	Scheme     string `mapstructure:"scheme"`
	Datacenter string `mapstructure:"datacenter"`
	Address    string `mapstructure:"address"`
}
type MetricServer struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
	Path string `mapstructure:"path"`
}
