package environment

import (
	"log"

	"github.com/spf13/viper"
)

func LoadEnvironment() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	viper.AutomaticEnv()
}

func GetString(key string) string {
	if !viper.IsSet(key) {
		panic("failed to get environment key: " + key)
	}

	return viper.GetString(key)
}

func GetInt(key string) int {
	if !viper.IsSet(key) {
		panic("failed to get environment key: " + key)
	}

	return viper.GetInt(key)
}

func GetBool(key string) bool {
	if !viper.IsSet(key) {
		panic("failed to get environment key: " + key)
	}

	return viper.GetBool(key)
}
