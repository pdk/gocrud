package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pdk/gocrud"

	_ "github.com/lib/pq"
)

const (
// user     = "***"
// password = "***"
// host     = "***"
// port     = 5432
// dbname   = "***"
// sslmode  = "***"
)

type Foo struct {
	ID      int64
	Name    string
	Age     int
	Address string
	Salary  float32
}

var (
	inserter1 = gocrud.NewAutoIncrIDInserter("foo", "ID", Foo{})
	inserter2 = gocrud.NewInserter("foo", Foo{})
	inserter3 func(*sql.DB, *Foo) (*Foo, error)
)

func init() {
	gocrud.MakeInserter(&inserter3, "foo", Foo{})
}

func (f Foo) Insert1(db *sql.DB) (Foo, error) {
	err := inserter1(db, &f)
	return f, err
}

func (f Foo) Insert2(db *sql.DB) (Foo, error) {
	err := inserter2(db, &f)
	return f, err
}

func (f Foo) Insert3(db *sql.DB) (*Foo, error) {
	return inserter3(db, &f)
}

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

	_, err = db.Exec("create table foo (id bigserial primary key, name text, age int, address varchar(200), salary numeric(12,2))")
	if err != nil {
		log.Fatalf("cannot create table: %v", err)
	}

	log.Printf("table created!")

	f := Foo{
		Name:    "Walter",
		Age:     32,
		Address: "123 Main St",
		Salary:  1234.56,
	}

	f, err = f.Insert1(db)
	if err != nil {
		log.Fatalf("failed to insert: %v", err)
	}

	log.Printf("rec inserted, new id = %d", f.ID)

	f.ID = 2
	f.Name = "Marty"
	f, err = f.Insert2(db)
	if err != nil {
		log.Fatalf("failed to insert again: %v", err)
	}

	log.Printf("2nd rec inserted")

	f.ID = 3
	f.Name = "Wart"
	nf, err := f.Insert3(db)
	if err != nil {
		log.Fatalf("failed to insert 3rd: %v", err)
	}

	log.Printf("final: %v", nf)
}
