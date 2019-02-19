package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var schema_person = `
CREATE TABLE person (
	id int(11) NOT NULL,
	first_name varchar(45) DEFAULT NULL,
	last_name varchar(45) DEFAULT NULL,
	email varchar(45) DEFAULT NULL,
	PRIMARY KEY (id)
);
`
var schema_palce = `
CREATE TABLE place (
	country varchar(45) DEFAULT NULL,
	city varchar(45) DEFAULT NULL,
	telcode int(11)
);
`

type Person struct {
	Id        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type Place struct {
	Country string         `db:"country"`
	City    sql.NullString `db:"city"`
	TelCode int            `db:"telcode"`
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("recover %s \n", r)
			os.Exit(1)
		}
	}()

	db, err := sqlx.Connect("mysql", "root:shzy2012@passwd@(localhost:3306)/test")
	if err != nil {
		panic(err.Error())
		//fatal will call os.Exit(1)
		//log.Fatalln(err)
	}

	db.MustExec(schema_person)
	db.MustExec(schema_palce)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (id,first_name, last_name, email) VALUES (?,?,?,?)", 1, "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (id,first_name, last_name, email) VALUES (?,?,?,?)", 2, "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES (?,?,?)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES (?,?)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES (?,?)", "Singapore", "65")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	tx.NamedExec("INSERT INTO person (id,first_name, last_name, email) VALUES (:id,:first_name, :last_name, :email)", &Person{3, "Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	// Query the database, storing results in a []Person (wrapped in []interface{})
	people := []Person{}
	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	log.Println("people ...")
	log.Println(people)
}
