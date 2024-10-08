package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local" env-required:"true"`
	DefaultStorage string     `yaml:"default_storage" env-default:"in-memory" env-required:"true"`
	HTTPServer     HTTPServer `yaml:"http_server"`
	GRPCServer     GRPCServer `yaml:"grpc_server"`
	Database       Database   `yaml:"database"`
	RabbitMQ       RabbitMQ   `yaml:"rabbitmq"`
	Scheduler      Scheduler  `yaml:"scheduler"`
}

type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type GRPCServer struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"50051"`
}

type RabbitMQ struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5672"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

type Scheduler struct {
	LaunchFrequency time.Duration `yaml:"launch_frequency" env-default:"1m"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	dir, _ := os.Getwd()
	fmt.Println(dir)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config file: ", err)
	}

	return &cfg
}
