package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

const (
	configPath = "./configs"
	configName = "main"
	envFile    = "./.env"
	Production = "production"
	Debug      = "debug"
)

type HTTP struct {
	Host            string
	Port            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func (h *HTTP) GetAddr() string {
	return fmt.Sprintf("%s:%s", h.Host, h.Port)
}

func (h *HTTP) GetReadTimeout() time.Duration {
	return h.ReadTimeout
}

func (h *HTTP) GetWriteTimeout() time.Duration {
	return h.WriteTimeout
}

func (h *HTTP) GetShutdownTimeout() time.Duration {
	return h.ShutdownTimeout
}

type DB struct {
	User            string
	Password        string
	Host            string
	Port            string
	Name            string
	Sslmode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

func (db *DB) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", db.User, db.Password, db.Host, db.Port, db.Name, db.Sslmode)
}

func (db *DB) GetMaxOpenCons() int {
	return db.MaxOpenConns
}

func (db *DB) GetMaxIdleCons() int {
	return db.MaxIdleConns
}

func (db *DB) GetConnMaxIdleTime() time.Duration {
	return db.ConnMaxIdleTime
}

func (db *DB) GetConnMaxLifetime() time.Duration {
	return db.ConnMaxLifetime
}

type Auth struct {
	Salt            string
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func (a *Auth) GetSalt() string {
	return a.Salt
}

func (a *Auth) GetSecret() string {
	return a.Secret
}

func (a *Auth) GetAccessTokenTTL() time.Duration {
	return a.AccessTokenTTL
}
func (a *Auth) GetRefreshTokenTTL() time.Duration {
	return a.RefreshTokenTTL
}

type Project struct {
	Name  string
	Level string
	Mode  string
}

func (p *Project) GetMode() string {
	return p.Mode
}

type Telemetry struct {
	GraylogPath string
}

func (t *Telemetry) GetGraylogPath() string {
	return t.GraylogPath
}

type Config struct {
	HTTP      HTTP
	DB        DB
	Auth      Auth
	Project   Project
	Telemetry Telemetry
}

var config = new(Config)
var once sync.Once

func init() {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("loading env file")
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("reading config err")
	}
}

func Get() *Config {
	once.Do(func() {
		var err error

		err = viper.Unmarshal(config)
		if err != nil {
			log.Fatal("reading config")
		}

		err = envconfig.Process("db", &config.DB)
		if err != nil {
			log.Fatal("error: get env for db")
		}

		err = envconfig.Process("auth", &config.Auth)
		if err != nil {
			log.Fatal("error: get env for auth")
		}

		config.Project.Mode = os.Getenv("MODE")
		if config.Project.Mode == "" {
			config.Project.Mode = Debug
		}
	})
	return config
}
