package config

import "github.com/kelseyhightower/envconfig"

type specification struct {
	Env            string   `default:"development" required:"true"`
	ServerHost     string   `default:"0.0.0.0" required:"true" split_words:"true"`
	ServerPort     int      `default:"8080" required:"true" split_words:"true"`
	AllowedOrigins []string `required:"true" split_words:"true"`
	DatabaseSource string   `required:"true" split_words:"true"`
	CacheAddress   string   `required:"true" split_words:"true"`
}

func LoadEnvVars() (*specification, error) {
	var spec specification
	if err := envconfig.Process("", spec); err != nil {
		return nil, err
	}
	return &spec, nil
}
