package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)

	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	return nil

}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
