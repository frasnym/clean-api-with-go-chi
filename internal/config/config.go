package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Constants struct {
	PORT       string
	PostgreSQL struct {
		URL      string
		DBName   string
		User     string
		Password string
	}
}

type Config struct {
	Constants
	Database *sql.DB
}

func New() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Constants = constants
	if err != nil {
		return &config, err
	}
	dbSession, err := sql.Open("postgres", "postgres://"+config.Constants.PostgreSQL.User+":"+config.Constants.PostgreSQL.Password+"@"+config.Constants.PostgreSQL.URL+"/"+config.Constants.PostgreSQL.DBName+"?sslmode=disable")
	if err != nil {
		return &config, err
	}
	if err = dbSession.Ping(); err != nil {
		panic(err)
	}
	config.Database = dbSession
	return &config, err
}

func initViper() (Constants, error) {
	viper.SetConfigName("todo.config") // Configuration fileName without the .TOML or .YAML extension
	viper.AddConfigPath(".")           // Search the root directory for the configuration file
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		return Constants{}, err
	}
	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.SetDefault("PORT", "8080")
	if err = viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	var constants Constants
	err = viper.Unmarshal(&constants)
	return constants, err
}
