package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	SMTP
	Code
	Cache
	Postgres
	Logger
	HTTPServer
}

type SMTP struct {
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
	Hostname string `env:"SMTP_HOSTNAME"`
	Port     string `env:"SMTP_PORT"`
}

type Code struct {
	Runes  string        `env:"CODE_RUNES"`
	Length int           `env:"CODE_LENGTH"`
	TTL    time.Duration `env:"CODE_DURATION"`
}

type Cache struct {
	CleanupTimeout time.Duration `env:"CACHE_CLEANUP_TIMEOUT"`
}

type Postgres struct {
	Username   string        `env:"POSTGRES_USER"`
	Password   string        `env:"POSTGRES_PASSWORD"`
	Host       string        `env:"POSTGRES_HOST"`
	Port       string        `env:"POSTGRES_PORT"`
	DB         string        `env:"POSTGRES_DB"`
	Timeout    time.Duration `env:"POSTGRES_TIMEOUT"`
	Migrations string        `env:"POSTGRES_MIGRATIONS"`
}

type Logger struct {
	InfoPath  string `env:"LOGGER_INFO_PATH"`
	ErrorPath string `env:"LOGGER_ERROR_PATH"`
}

type HTTPServer struct {
	Interface       string        `env:"SERVER_INTERFACE"`
	Port            string        `env:"SERVER_PORT"`
	ReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT"`
}

func New() (*Config, error) {
	conf := &Config{}

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(conf); err != nil {
		return nil, err
	}

	return conf, nil
}