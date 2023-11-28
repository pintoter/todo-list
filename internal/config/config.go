package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	configPath = "./config/config.yml"
)

type HTTP struct {
	Host            string
	Port            string
	ShutdownTimeout time.Duration
}

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Sslmode  string
}

func (db *DB) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", db.User, db.Password, db.Host, db.Port, db.Name, db.Sslmode)
}

type Config struct {
	HTTP HTTP
	DB   DB
}

var config = new(Config)
var once sync.Once

func init() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("reading config err")
	}
}

func Get() *Config {
	once.Do(func() {

		err := viper.Unmarshal(config)
		if err != nil {
			log.Fatal("reading config")
		}
	})
	return config
}
