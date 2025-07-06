package config

import (
	"github.com/spf13/viper"

	"urlshort/internal/utils"
)

type Config struct {
	Database  string
	JwtSecret []byte
	Redis     Redis
	Port      string
}

type Redis struct {
	Enabled  bool
	Host     string
	Port     int
	Password string
	DB       int
	TTL      int
}

func LoadConfig() Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	utils.CheckError(err)

	return Config{
		Database:  getDatabase(),
		JwtSecret: getJWTSecret(),
		Redis:     GetRedis(),
		Port:      viper.GetString("server.port"),
	}
}

func GetRedis() Redis {
	enabled := viper.GetBool("redis.enabled")
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	ttl := viper.GetInt("redis.TTL")

	return Redis{
		Enabled:  enabled,
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
		TTL:      ttl,
	}
}

func getJWTSecret() []byte {
	return []byte(viper.GetString("jwt.secret"))
}
