package db

import (
	"api-server/config"
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SqlDB *sql.DB

func init() {
	if DB == nil {
		SqlDB, err := sql.Open("postgres", config.Get().DatabaseURL)
		if err != nil {
			log.Fatal("Unable to open postges connection. Err:", err)
		}

		SqlDB.SetMaxIdleConns(5)
		SqlDB.SetMaxOpenConns(5)
		SqlDB.SetConnMaxLifetime(time.Hour)

		DB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: SqlDB,
		}), &gorm.Config{})
		if err != nil {
			log.Fatal("Unable to open postges gorm connection. Err:", err)
		}

		log.Println("Successfully established database connection")
	}
}

type DBConn struct {
	*gorm.DB
}

func New() *DBConn {
	return &DBConn{
		DB: DB,
	}
}
