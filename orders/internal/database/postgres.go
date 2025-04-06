package database

import (
	"fmt"
	"github.com/TeslaMode1X/firstAssignmentADVPROG/orders/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

type postgresDatabase struct {
	Db *gorm.DB
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return p.Db
}

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			conf.DB.Host,
			conf.DB.User,
			conf.DB.Password,
			conf.DB.DatabaseName,
			conf.DB.Port,
			conf.DB.SSLMode,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("failed to connect database: %v", err))
		}

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}
