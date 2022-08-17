package config

import "github.com/spf13/viper"

type Configuration struct {
	ProgramId string
	Endpoint  string
}

func LoadConfig(configPath string) Configuration {
	conf := Configuration{}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return conf
}
