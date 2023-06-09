package main

import (
	"database/sql"
	"time"
	"os"

	"github.com/go-sql-driver/mysql"
)

// ConnectDB connects to a MySQL database, given the following
// environment variables are set: DBHOST, DBPORT, DBUSER, DBPASS,
// DBNAME
func ConnectDB() (*sql.DB, error) {
	addr := os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT")
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	return db, nil
}
