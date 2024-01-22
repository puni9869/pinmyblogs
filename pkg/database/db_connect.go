package database

import (
	"fmt"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func NewDBLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	return newLogger
}

func NewConnection(cfg *config.Database) (*gorm.DB, error) {
	//follow the below connectionString
	// "type://username:password@host:port/dbName"
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
		dbLogger = NewDBLogger()
	}
	ormConfig := gorm.Config{
		Logger: dbLogger,
	}
	db, err := gorm.Open(postgres.Open(dsn), &ormConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}
