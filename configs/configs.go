package configs

import "os"

type Config struct {
	DB  DBConfig
	App AppConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type AppConfig struct {
	Port string
	Host string
}

func LoadConfig() *Config {
	cfg := &Config{}
	cfg.App.Host = os.Getenv("APP_HOST")
	cfg.App.Port = os.Getenv("APP_PORT")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Username = os.Getenv("DB_USERNAME")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.DBName = os.Getenv("DB_NAME")
	return cfg
}