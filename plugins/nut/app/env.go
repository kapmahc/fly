package app

import "github.com/spf13/viper"

// Name get server.name
func Name() string {
	return viper.GetString("server.name")
}

// IsProduction production mode ?
func IsProduction() bool {
	return viper.GetString("env") == "production"
}
