package config

var StartConfig *Config

func Init(configFile string) error {
	var err error

	// copy default config
	StartConfig = new(Config)
	*StartConfig = *DefaultConfig

	if configFile != "" {
		err = FromFile(StartConfig, configFile)
	} else {
		err = FromFileDefault(StartConfig)
	}
	FromFlags(StartConfig)

	return err
}
