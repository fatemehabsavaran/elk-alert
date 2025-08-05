package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

func Load[T any]() *T {
	viper.SetConfigName(GetEnv("CONFIG_NAME", "config"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	t := new(T)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	} else {
		log.Printf("Config file loaded successfully")
	}

	if err := viper.Unmarshal(&t); err != nil {
		log.Fatalf("Error unmarshalling config: %s", err)
	} else {
		log.Printf("Config file unmarshalled successfully")
	}

	return t
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
