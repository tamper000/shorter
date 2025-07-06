//go:build !postgre

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func getDatabase() string {
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.db")

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
}
