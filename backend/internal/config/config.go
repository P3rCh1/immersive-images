package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var Config struct {
	Logger Logger `yaml:"logger"`
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
	S3     S3     `yaml:"s3"`
	Images Images `yaml:"images"`
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

	Password string `yaml:"-" env:"POSTGRES_PASSWORD" validate:"required"`
}

type S3 struct {
	Endpoint string `yaml:"endpoint" env-default:"https://s3.ru-3.storage.selcloud.ru"`
	Port     string `yaml:"port"     env-default:"443"`
	Region   string `yaml:"region"   env-default:"ru-3"`
	Bucket   string `yaml:"bucket"   env-default:"immersive-images"`

	AccessKey string `yaml:"-" env:"AWS_ACCESS_KEY" validate:"required"`
	SecretKey string `yaml:"-" env:"AWS_SECRET_KEY" validate:"required"`
}

type Images struct {
	DefaultLimit int `yaml:"default_limit" env-default:"10"`
	MaxLimit     int `yaml:"max_limit"     env-default:"100"`

	DefaultStyle   string `yaml:"default_style"   env-default:"PATCHES"`
	DefaultPalette string `yaml:"default_palette" env-default:"NEON"`

	DefaultWidth  int `yaml:"default_width"  env-default:"256"`
	DefaultHeight int `yaml:"default_height" env-default:"256"`
	MinWidth      int `yaml:"min_width"      env-default:"64"`
	MaxWidth      int `yaml:"max_width"      env-default:"1024"`
	MinHeight     int `yaml:"min_height"     env-default:"64"`
	MaxHeight     int `yaml:"max_height"     env-default:"1024"`

	DefaultScale float64 `yaml:"default_scale" env-default:"3"`
	MinScale     float64 `yaml:"min_scale"     env-default:"0.5"`
	MaxScale     float64 `yaml:"max_scale"     env-default:"10"`
}
