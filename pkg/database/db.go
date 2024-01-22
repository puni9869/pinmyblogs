package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// dbObj is the instance of the database
var dbObj *gorm.DB

// NewLogger returns the SQL logger configuration
func NewLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Threshold for slow SQL queries
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound errors in the logger
			ParameterizedQueries:      true,        // Exclude parameters from the SQL log
			Colorful:                  false,       // Disable color in logs
		},
	)
	return newLogger
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		cfg.Type,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	var dbLogger logger.Interface

	if cfg.LogSql {
		dbLogger = NewLogger()
	}

	ormConfig := gorm.Config{
		Logger: dbLogger,
	}

	db, err := gorm.Open(postgres.Open(dsn), &ormConfig)
	if err != nil {
		return nil, err
	}

	dbObj = db
	return db, nil
}

// RegisterModels configures the available models
func RegisterModels(db *gorm.DB) *gorm.DB {
	_ = db.AutoMigrate(&models.User{})
	return db
}

// Db returns the global database instance
func Db() *gorm.DB {
	return dbObj
}
