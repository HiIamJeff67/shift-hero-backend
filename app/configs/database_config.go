package configs

import util "github.com/your-org/go-start-monolithic-kit/app/util"

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string // the port inside the container, so please leave this as 5432 for PostgreSQL
}

var (
	PostgresDatabaseConfig = DatabaseConfig{
		Host:     util.GetEnv("DB_HOST", "go-start-monolithic-kit-db"),
		User:     util.GetEnv("DB_USER", "master"),
		Password: util.GetEnv("DB_PASSWORD", ""),
		DBName:   util.GetEnv("DB_NAME", "go-start-monolithic-kit-db"),
		Port:     util.GetEnv("DOCKER_DB_PORT", "5432"),
	}
)
