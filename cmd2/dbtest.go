package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	user     = "***"
	password = "***"
	host     = "***"
	port     = 0000
	dbname   = "***"
	sslmode  = "***"
)

func main() {

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("drop table if exists foo")
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	_, err = db.Exec("create table foo (id bigserial primary key, name text, age int, address varchar(200), salary numeric(12,0))")
	if err != nil {
		log.Fatalf("cannot create table: %v", err)
	}

}
