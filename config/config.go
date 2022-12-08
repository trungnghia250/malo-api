package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Schema struct {
	Mongodb struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
		User     string `mapstructure:"user"`
		Pass     string `mapstructure:"pass"`
	} `mapstructure:"mongodb"`

	Encryption struct {
		OIDKey          string `mapstructure:"oid_key"`
		SigningMethod   string `mapstructure:"signning_method"`
		JWTSecret       string `mapstructure:"jwt_secret"`
		JWTExp          int    `mapstructure:"jwt_exp"`
		JWTPol          string `mapstructure:"jwt_pol"`
		RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
	} `mapstructure:"encryption"`

	RabbitMq RabbitMQ `mapstructure:"rabbit_instance"`
}

type RabbitMQ struct {
	AMQP         string `mapstructure:"amqp"`
	Exchange     string `mapstructure:"exchange"`
	RoutingKey   string `mapstructure:"routing_key"`
	Queue        string `mapstructure:"queue"`
	Tags         string `mapstructure:"tag"`
	TimeOutRetry int    `mapstructure:"timeout_retry"`
}

var Config Schema

func NewSchema() *Schema {
	schema := new(Schema)
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")       // Look for config in current directory
	config.AddConfigPath("config/") // Optionally look for config in the working directory.

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&schema)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return schema
}

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")             // Look for config in current directory
	config.AddConfigPath("config/")       // Optionally look for config in the working directory.
	config.AddConfigPath("../config/")    // Look for config needed for tests.
	config.AddConfigPath("../")           // Look for config needed for tests.
	config.AddConfigPath("../../config/") // used for integration_test

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&Config)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
