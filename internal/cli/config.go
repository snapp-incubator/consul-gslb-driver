package cli

// corespoding struct for config.yaml file
type Config struct {
	Verbose int    `mapstructure:"verbose"`
	Listen  Listen `mapstructure:"listen"`
}

type Listen struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

// Creates a Config struct with setting default values.
// If a default pflag value is set in cobra, this is being overrided
// so those variables MUST NOT be declared here.
func newConfig() *Config {
	return &Config{
		Listen: Listen{
			IP: "127.0.0.1",
		},
	}
}
