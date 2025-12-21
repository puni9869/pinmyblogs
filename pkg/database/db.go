package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"

	sqliteGo "github.com/mattn/go-sqlite3"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"
)

var (
	// dbObj is the instance of the database
	dbObj *gorm.DB
	once  sync.Once
)

// newLogger returns the SQL logger configuration
func newLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Threshold for slow SQL queries
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound errors in the logger
			ParameterizedQueries:      true,        // Exclude parameters from the SQL log
			Colorful:                  false,       // Disable color in logs
		},
	)
}

// NewPostgresConnection creates a new database postgres connection
func NewPostgresConnection(cfg *config.DatabaseObj) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		cfg.Type,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName)
	var dbLogger logger.Interface
	if cfg.LogSql {
		dbLogger = newLogger()
	}
	ormConfig := gorm.Config{
		Logger: dbLogger,
	}

	db, err := gorm.Open(postgres.Open(dsn), &ormConfig)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}
	dbObj = db
	// making an object singleton
	once.Do(func() { dbObj = db })
	return dbObj, nil
}

// NewSqliteConnection creates a new database sqlite connection
func NewSqliteConnection(cfg *config.DatabaseObj) (*gorm.DB, error) {
	ormConfig := gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	if cfg.LogSql {
		ormConfig.Logger = newLogger()
	}

	const CustomDriverName = "sqlite3_extended"
	sql.Register(CustomDriverName,
		&sqliteGo.SQLiteDriver{
			ConnectHook: func(conn *sqliteGo.SQLiteConn) error {
				err := conn.RegisterFunc(
					"gen_random_uuid",
					func(arguments ...interface{}) (string, error) {
						return uuid.New().String(), nil // Return a string value.
					},
					true,
				)
				return fmt.Errorf("open database connection: %w", err)
			},
		},
	)
	conn, err := sql.Open(CustomDriverName, cfg.FileName)
	if err != nil {
		return nil, fmt.Errorf("sql open failed: %w", err)
	}

	db, err := gorm.Open(sqlite.Dialector{
		DriverName: CustomDriverName,
		DSN:        fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL", cfg.FileName),
		Conn:       conn,
	}, &ormConfig)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}
	// making an object singleton
	once.Do(func() { dbObj = db })
	return dbObj, nil
}

// RegisterModels configures the available models
func RegisterModels(db *gorm.DB) {
	// m is list of all the database models
	m := []any{
		&models.User{},
		&models.Session{},
		&models.Url{},
		&models.Setting{},
		&models.JoinWaitList{},
	}
	if err := db.AutoMigrate(m...); err != nil {
		panic(err)
	}
}

// Db returns the global database instance
func Db() *gorm.DB {
	return dbObj
}
