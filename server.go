// (c) 2016 Written by Chris J Arges <christopherarges@gmail.com>

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type Business struct {
	Id         int
	Uuid       string
	Name       string
	Address    string
	Address2   string
	City       string
	Zip        string
	Country    string
	Phone      string
	Website    string
	Created_at string
}

type Meta struct {
	Start int	// starting id of results
	Length int	// length of results returned
	Total int	// total results in the database
}

var db *sql.DB

// GetBusiness returns a single business based on an id
func GetBusinesses(id int) (b Business, err error) {
	rows, err := db.Query("SELECT id,uuid,name,address,address2,city FROM businesses WHERE id=" + strconv.Itoa(id) + ";")
	if err != nil {
		return
	}
	rows.Next()
	err = rows.Scan(&b.Id, &b.Uuid, &b.Name, &b.Address, &b.Address2, &b.City)
	return
}

// BusinessesHandler is the business end of the businesses handler.
func BusinessesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if len(id) > 0 {
		i, _ := strconv.Atoi(id)
		b, _ := GetBusinesses(i)
		fmt.Fprintf(w, "%+v\n", b)
		return
	}
}

func main() {
	// Load database
	var err error
	db, err = sql.Open("sqlite3", "businesses.sql")
	if err != nil {
		panic(err)
	}

	// Serve API
	router := httprouter.New()
	router.GET("/businesses", BusinessesHandler)
	router.GET("/businesses/:id", BusinessesHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
