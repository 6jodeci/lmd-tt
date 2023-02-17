package cofig

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBHost string `env:"DB_HOST" env-default:"localhost"`
	DBPort int    `env:"DB_PORT" env-default:"5432"`
	DBUser string `env:"DB_USER"`
	DBPass string `env:"DB_PASS"`
	DBName string `env:"DB_NAME"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("../app.env", &instance); err != nil {
			log.Fatal("failed to read config", err)
		}
	})
	if instance == nil {
		log.Fatal("failed to initialize config")
	}
	return instance
}
