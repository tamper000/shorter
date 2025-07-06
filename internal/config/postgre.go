//go:build postgre

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func getDatabase() string {
	host := viper.GetString("postgre.host")
	port := viper.GetInt("postgre.port")
	user := viper.GetString("postgre.user")
	password := viper.GetString("postgre.password")
	database := viper.GetString("postgre.db")

	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=require", user, password, host, port, database)
}
