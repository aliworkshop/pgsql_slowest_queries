package application

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type SqlConfig struct {
	Dialect            string
	Host               string
	Port               string
	Username           string
	Password           string
	DbName             string
	ConnectionString   string
	MaxIdleConnections *int
	MaxOpenConnections *int
	MaxLifetimeSeconds *int
}

func NewPostgreSqlConnection(config SqlConfig) *gorm.DB {
	connStr := config.ConnectionString
	if connStr == "" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			config.Host,
			config.Port,
			config.Username,
			config.DbName,
			config.Password,
		)
	}
	dialect := postgres.New(postgres.Config{
		DSN:                  connStr,
		PreferSimpleProtocol: true,
	})
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             5 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn,     // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)
	d, err := gorm.Open(dialect, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := d.DB()
	if err != nil {
		panic(err)
	}
	if config.MaxIdleConnections != nil {
		sqlDB.SetMaxIdleConns(*config.MaxIdleConnections)
	}
	if config.MaxOpenConnections != nil {
		sqlDB.SetMaxOpenConns(*config.MaxOpenConnections)
	}
	if config.MaxLifetimeSeconds != nil {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(*config.MaxLifetimeSeconds))
	}
	return d.Debug()
}
