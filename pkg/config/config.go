package config

import (
	"bytes"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Address  string `yaml:"address"`
	Iface    string `yaml:"iface"`
	LogLevel string `yaml:"loglevel"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Address:  "0.0.0.0:8080",
		Iface:    "",
		LogLevel: "info",
	}
}

func LoadConfig() (config Config, err error) {
	b, err := yaml.Marshal(NewDefaultConfig())
	if err != nil {
		return
	}
	defaultConfig := bytes.NewReader(b)

	viper.SetConfigType("yaml")
	err = viper.MergeConfig(defaultConfig)
	if err != nil {
		return
	}

	viper.SetConfigName("xdpdropper")
	viper.SetConfigFile(".")
	err = viper.MergeInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigParseError)
		if ok {
			return
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("XDPDROPPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	err = viper.Unmarshal(&config)
	return
}
