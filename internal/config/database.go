package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	DATABASE_URL           = "DATABASE_URL"
	MIGRATION_PATH         = "MIGRATION_PATH"
	DATABASE_SSL_MODE      = "POSTGRES_SSLMODE"
	DATABASE_HOST          = "POSTGRES_HOST"
	DATABASE_PORT          = "POSTGRES_PORT"
	DATABASE_NAME          = "POSTGRES_DB"
	DATABASE_PASSWORD      = "POSTGRES_PASSWORD"
	DATABASE_USERNAME      = "POSTGRES_USER"
	DATABASE_PASSWORD_FILE = "POSTGRES_PASSWORD_FILE"
)

// the database configuration struct: defines the database connection parameters
type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     uint16
	Name     string
	SSLMode  string
}

func (c *DatabaseConfig) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func (c *DatabaseConfig) Validate() error {
	if c.Username == "" {
		return fmt.Errorf("database username is required")
	}

	if c.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Port == 0 {
		return fmt.Errorf("database port is required")
	}

	if c.Name == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
}

// load the database password from the environment
// loads the password from the DATABASE_PASSWORD environment variable if it is set
// if not, it loads the password from the DATABASE_PASSWORD_FILE environment variable
func loadDatabasePassword() (string, error) {
	password, ok := os.LookupEnv(DATABASE_PASSWORD)
	if ok {
		return password, nil
	}

	passwordFile, ok := os.LookupEnv(DATABASE_PASSWORD_FILE)
	if !ok {
		return "", fmt.Errorf("failed to load database password: %s and %s env var aren't set", DATABASE_PASSWORD, DATABASE_PASSWORD_FILE)
	}

	data, err := os.ReadFile(passwordFile)

	if err != nil {
		return "", fmt.Errorf("failed to read password file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

// creates a new database configuration struct based on the environment variables
// returns an error if any of the required environment variables are missing
func NewDatabaseConfig() (*DatabaseConfig, error) {
	password, err := loadDatabasePassword()
	if err != nil {
		return nil, err
	}

	username, ok := os.LookupEnv(DATABASE_USERNAME)
	if !ok {
		return nil, fmt.Errorf("failed to load database username: %s env var isn't set", DATABASE_USERNAME)
	}

	host, ok := os.LookupEnv(DATABASE_HOST)
	if !ok {
		return nil, fmt.Errorf("failed to load database host: %s env var isn't set", DATABASE_HOST)
	}

	portStr, ok := os.LookupEnv(DATABASE_PORT)
	if !ok {
		return nil, fmt.Errorf("failed to load database port: %s env var isn't set", DATABASE_PORT)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database port: %w", err)
	}

	name, ok := os.LookupEnv(DATABASE_NAME)
	if !ok {
		return nil, fmt.Errorf("failed to load database name: %s env var isn't set", DATABASE_NAME)
	}

	sslMode, ok := os.LookupEnv(DATABASE_SSL_MODE)
	if !ok {
		sslMode = "disable"
	}

	config := &DatabaseConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     uint16(port),
		Name:     name,
		SSLMode:  sslMode,
	}

	return config, config.Validate()
}
