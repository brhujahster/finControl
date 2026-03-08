package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("sqlite", "./fincontrol.db")
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Banco de dados inacessível: %v", err)
	}

	log.Println("Banco de dados conectado.")
}
