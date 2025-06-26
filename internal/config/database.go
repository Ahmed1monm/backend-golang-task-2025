package config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/errors"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	Env          string        // Environment (development/production)
	
	// Connection pool settings
	MaxIdleConns int           // Maximum number of idle connections
	MaxOpenConns int           // Maximum number of open connections
	MaxIdleTime  time.Duration // Maximum amount of time a connection may be idle
	MaxLifetime  time.Duration // Maximum amount of time a connection may be reused
	
	// Connection retry settings
	MaxRetries   int           // Maximum number of connection retries
	RetryBackoff time.Duration // Initial backoff duration between retries
}

// NewDBConfig creates a new database configuration
func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:         utils.GetEnv("DB_HOST", "localhost"),
		Port:         utils.GetEnv("DB_PORT", "5432"),
		User:         utils.GetEnv("DB_USER", "postgres"),
		Password:     utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:       utils.GetEnv("DB_NAME", "myapp"),
		SSLMode:      utils.GetEnv("DB_SSLMODE", "disable"),
		Env:          utils.GetEnv("ENV", "development"),

		// Connection pool settings
		MaxIdleConns: utils.GetEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns: utils.GetEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		MaxIdleTime:  utils.GetEnvAsDuration("DB_MAX_IDLE_TIME", 30*time.Minute),
		MaxLifetime:  utils.GetEnvAsDuration("DB_MAX_LIFETIME", 1*time.Hour),
		
		// Connection retry settings
		MaxRetries:   utils.GetEnvAsInt("DB_MAX_RETRIES", 5),
		RetryBackoff: utils.GetEnvAsDuration("DB_RETRY_BACKOFF", 1*time.Second),
	}
}

// Connect establishes a connection to the database with retry mechanism
func (c *DBConfig) Connect() (*gorm.DB, error) {
	ctx := context.Background()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)

	var db *gorm.DB
	var err error
	backoff := c.RetryBackoff

	for attempt := 1; attempt <= c.MaxRetries; attempt++ {
		logger.Info(ctx, "Attempting database connection",
			zap.Int("attempt", attempt),
			zap.Int("max_retries", c.MaxRetries))

		// Configure GORM
		config := &gorm.Config{}
		if c.Env == "development" {
			config.Logger = logger.NewGormLogger()
		}

		// Attempt to connect
		db, err = gorm.Open(postgres.Open(dsn), config)
		if err == nil {
			// Successfully connected
			break
		}

		if attempt == c.MaxRetries {
			return nil, errors.NewServerError(
				fmt.Sprintf("failed to connect to database after %d attempts", c.MaxRetries),
				err,
				http.StatusInternalServerError,
			)
		}

		logger.Warn(ctx, "Database connection failed, retrying",
			zap.Error(err),
			zap.Duration("backoff", backoff))

		// Wait with exponential backoff before next attempt
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.NewServerError("error getting database instance", err, http.StatusInternalServerError)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(c.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(c.MaxLifetime)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, errors.NewServerError("error verifying database connection", err, http.StatusInternalServerError)
	}

	logger.Info(ctx, "Successfully connected to database",
		zap.String("host", c.Host),
		zap.String("database", c.DBName))

	return db, nil
}
