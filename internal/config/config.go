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
	Redis
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
	CacheType      string        `env:"CACHE_TYPE"`
	CleanupTimeout time.Duration `env:"CACHE_CLEANUP_TIMEOUT"`
	Redis
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

type Redis struct {
	DB       int    `env:"REDIS_DB"`
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Username string `env:"REDIS_USER"`
	Password string `env:"REDIS_PASSWORD"`
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
