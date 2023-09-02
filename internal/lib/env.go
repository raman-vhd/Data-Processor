package lib

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerPort           string `mapstructure:"SERVER_PORT"`
	KafkaBootstrapServer string `mapstructure:"KAFKA_BOOTSTRAP_SERVER"`
}

func NewEnv() Env {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed reading env variables: %v\n", err)
	}

	var env Env
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("failed unmarshaling env variables: %v\n", err)
	}

	return env
}
