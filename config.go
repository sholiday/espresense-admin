package eadmin

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Broker MqttBroker
}

type MqttBroker struct {
	Server   string // Like tcp://host:port.
	Username string
	Password string
	ClientID string
}

type ServerConfig struct {
	Host string
	Port int
}

func LoadConfig() (Config, error) {
	viper.SetDefault("server.port", 12312)
	viper.SetDefault("broker.prefix", "espresense")
	viper.SetEnvPrefix("ESPRESENSE_ADMIN")
	viper.SetConfigName("espresense-admin")
	viper.AddConfigPath("$HOME/.config/espresense-admin/")
	viper.AddConfigPath("/etc/espresense-admin/")
	viper.AddConfigPath("/config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("Failed to read config file: %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("Failed to unmarshal config file: %s", err)
	}
	return config, err
}
