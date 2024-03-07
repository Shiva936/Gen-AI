package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Config struct {
	AppName          string `json:"app_name"`
	Port             string `json:"port"`
	LogLevel         string `json:"log_level"`
	MessageTopic     string `json:"message_topic"`
	MessageBrokerURL string `json:"message_broker_url"`
	DatabaseURL      string `json:"database_url"`
}

var conf *Config

func Get() *Config {
	return conf
}

func Set(c *Config) {
	conf = c
	conf.MessageTopic = os.Getenv("MESSAGE_TOPIC")
	conf.MessageBrokerURL = os.Getenv("MESSAGE_BROKER_URL")
	conf.DatabaseURL = os.Getenv("DATABASE_URL")
}

func ParseJSON(file io.Reader, v any) error {
	data, err := io.ReadAll(file)
	if err != nil {
		return errors.New("unable to read file")
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return errors.New("unable to unmarshal")
	}

	return nil
}
