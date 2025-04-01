package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	// "github.com/joho/godotenv"
)

type DbInfo struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	MaxConn  int32  `yaml:"max_conn" env-default:"10"`
}

type HttpInfo struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type Config struct {
	Storage    string   `yaml:"storage" env-default:"inmemory"`
	HttpServer HttpInfo `yaml:"http_server"`
	DB         DbInfo   `yaml:"postgres"`
}

func LoadConfig() (*Config, error) {
	// err := godotenv.Load("/.env")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to load .env file: %s", err)
	// }

	var cfg Config
	if err := cleanenv.ReadConfig("config/main.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	log.Println("CFG", cfg)

	// cfg.DB.Password = os.Getenv("DB_PASSWORD")
	// cfg.DB.Host = os.Getenv("DB_HOST")
	// cfg.DB.Port = os.Getenv("DB_PORT")
	// cfg.DB.Name = os.Getenv("DB_NAME")
	// cfg.DB.User = os.Getenv("DB_USER")

	return &cfg, nil
}
