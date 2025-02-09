package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Pinger struct {
		PingInterval time.Duration `mapstructure:"ping-interval"`
	} `mapstructure:"pinger"`
	Broker struct {
		URL   string `mapstructure:"url"`
		Queue string `mapstructure:"queue"`
	} `mapstructure:"broker"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.BindEnv("pinger.ping-interval", "PING_INTERVAL")
	viper.BindEnv("broker.url", "BROKER_URL")
	viper.BindEnv("broker.queue", "BROKER_QUEUE")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}
	var config Config
	if err := viper.Unmarshal(&config, func(decoder *mapstructure.DecoderConfig) { decoder.ErrorUnset = true }); err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return config, nil
}
