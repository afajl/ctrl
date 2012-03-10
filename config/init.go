package config

var StartConfig *Config

func init() {
	// copy default config
	StartConfig = new(Config)
	*StartConfig = *DefaultConfig
}

func Init(configFile string) error {
	var err error

	if configFile != "" {
		err = FromFile(StartConfig, configFile)
	} else {
		err = FromFileDefault(StartConfig)
	}
	FromFlags(StartConfig)

	return err
}
