package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var Config struct {
	Logger Logger `yaml:"logger"`
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
}

func InitConfig(path string) error {
	return cleanenv.ReadConfig(path, &Config)
}

type Logger struct {
	Level  string `yaml:"level"  env-default:"warn"`
	Format string `yaml:"format" env-default:"json"`
}

type Server struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"8080"`

	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"20s"`

	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"30s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host string `yaml:"host" env-default:"postgres"`
	Port string `yaml:"port" env-default:"5432"`
	Name string `yaml:"name" env-default:"postgres"`

	Password string `yaml:"-" env:"PASSWORD" validate:"required"`
}
