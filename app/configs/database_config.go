package configs

import util "github.com/HiIamJeff67/shift-hero-backend/app/util"

type DatabaseConfig struct {
	URL      string
	Host     string
	User     string
	Password string
	DBName   string
	Port     string // the port inside the container, so please leave this as 5432 for PostgreSQL
}

var (
	PostgresDatabaseConfig = DatabaseConfig{
		URL:      util.GetEnv("DATABASE_URL", ""),
		Host:     util.GetEnv("DB_HOST", "shift-hero-db"),
		User:     util.GetEnv("DB_USER", "master"),
		Password: util.GetEnv("DB_PASSWORD", ""),
		DBName:   util.GetEnv("DB_NAME", "shift-hero-db"),
		Port:     util.GetEnv("DB_PORT", util.GetEnv("DOCKER_DB_PORT", "5432")),
	}
)
