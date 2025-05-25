package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AppGofemart AppGofemart    `validate:"required"`
	AppMigrator AppMigrator    `validate:"required"`
	Infra       Infrastructure `validate:"required"`
}

type AppGofemart struct {
	AppMode     string        `env:"APP__GOFEMART__MODE" validate:"required,oneof=dev prod local"`
	AppName     string        `env:"APP__GOFEMART__NAME" validate:"required,min=3"`
	LogLevel    string        `env:"APP__GOFEMART__LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	AppPort     string        `env:"APP__GOFEMART__PORT" validate:"required,numeric,min=4,max=5"`
	AppHost     string        `env:"APP__GOFEMART__HOST" validate:"required,hostname_rfc1123|ipv4|ipv6"`
	JWTSecret   string        `env:"APP__GOFEMART__JWT_SECRET" validate:"required"`
	JWTIssuer   string        `env:"APP__GOFEMART__JWT_ISSUER"    validate:"required,url"`
	JWTTokenTTL time.Duration `env:"APP__GOFEMART__JWT_TOKEN_TTL" validate:"required,gt=0"`
}

type AppMigrator struct {
	LogLevel       string `env:"APP__MIGRATOR__LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	AppMode        string `env:"APP__MIGRATOR__MODE" validate:"required,oneof=dev prod local"`
	ServiceName    string `env:"APP__MIGRATOR__NAME" validate:"required,min=3"`
	MigrationsPath string `env:"APP__MIGRATOR__MIGRATIONS_PATH" validate:"required"`
}

type Infrastructure struct {
	Postgres PostgreSQL `validate:"required"`
	Redis    Redis      `validate:"required"`
	Jaeger   Jaeger     `validate:"required"`
}

type PostgreSQL struct {
	Host     string `env:"INFRA__POSTGRES__HOST" validate:"required,hostname|ip"`
	Port     string `env:"INFRA__POSTGRES__PORT" validate:"required,numeric"`
	User     string `env:"INFRA__POSTGRES__USER" validate:"required"`
	Password string `env:"INFRA__POSTGRES__PASSWORD" validate:"required"`
	Database string `env:"INFRA__POSTGRES__DATABASE" validate:"required,alphaunicode"`
	ConnStr  string `env:"INFRA__POSTGRES__CONN_STR" validate:"required,url"`
}

type Redis struct {
	Host         string `env:"INFRA__REDIS__HOST" validate:"required,hostname|ip"`
	Port         string `env:"INFRA__REDIS__PORT" validate:"required,numeric"`
	Password     string `env:"INFRA__REDIS__PASSWORD"`
	Database     int    `env:"INFRA__REDIS__DATABASE" validate:"gte=0,lte=15"`
	PoolSize     int    `env:"INFRA__REDIS__POOL_SIZE" validate:"gte=1,lte=100"`
	MinIdleConns int    `env:"INFRA__REDIS__MIN_IDLE_CONNS" validate:"gte=0"`
	PoolTimeout  int    `env:"INFRA__REDIS__POOL_TIMEOUT" validate:"gte=1"`
}

type Jaeger struct {
	ServiceName string `env:"JAEGER_SERVICE_NAME" validate:"required"`
	Host        string `env:"JAEGER_AGENT_HOST" validate:"required,hostname|ip"`
	Port        string `env:"JAEGER_AGENT_PORT" validate:"required,numeric"`
	LogSpans    bool   `env:"JAEGER_LOG_SPANS"`
	Enabled     bool   `env:"JAEGER_ENABLED"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
