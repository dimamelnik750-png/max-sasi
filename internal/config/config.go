package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Port       string
	Env        string
	JWTSecret  string `mapstructure:"jwt_secret"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	RedisHost  string `mapstructure:"redis_host"`
	RedisPort  int    `mapstructure:"redis_port"`
}

func Load() Config {
	viper.SetConfigFile(findEnvFile())
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error reading .env file:", err)
	}

	return Config{
		Port:       viper.GetString("PORT"),
		Env:        viper.GetString("ENV"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		RedisHost:  viper.GetString("REDIS_HOST"),
		RedisPort:  viper.GetInt("REDIS_PORT"),
	}
}

func findEnvFile() string {
	candidates := []string{
		".env",
		filepath.Join("..", ".env"),
	}

	wd, err := os.Getwd()
	if err == nil {
		candidates = append(candidates,
			filepath.Join(wd, ".env"),
			filepath.Join(wd, "..", ".env"),
		)
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ".env"
}
