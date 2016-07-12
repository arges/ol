// (c) 2016 Written by Chris J Arges <christopherarges@gmail.com>

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Business struct {
	Id         int    `json:"id"`
	Uuid       string `json:"uuid"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	Zip        string `json:"zip"`
	Country    string `json:"country"`
	Phone      string `json:"phone"`
	Website    string `json:"website"`
	Created_at string `json:"created_at"`
}

type Meta struct {
	Page   int // which page of results has been returned
	Length int // length of results returned
}

type All struct {
	Businesses []Business `json:"businesses"` // All businesses returned
	Meta       Meta       `json:"meta"`       // Meta has page/length info
}

type Context struct {
	db     *sql.DB
	router *mux.Router
}

var SqlString string = "SELECT id,uuid,name,address,address2,city,zip,country,phone,website,created_at FROM businesses"

// FormatJson returns formatted json given an object
func FormatJson(object interface{}) (out bytes.Buffer) {
	b, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}
	json.Indent(&out, b, "", "\t")
	return
}

// GetBusiness returns a single business based on an id
func (c *Context) GetBusiness(id int) (b Business, err error) {
	rows, err := c.db.Query(SqlString + " WHERE id=" + strconv.Itoa(id) + ";")
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&b.Id, &b.Uuid, &b.Name, &b.Address, &b.Address2, &b.City, &b.Zip, &b.Country, &b.Phone, &b.Website, &b.Created_at)
	if err != nil {
		return
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// GetBusinesses returns a single business based on an id
func (c *Context) GetBusinesses(page int, length int) (bs []Business, err error) {
	start := page * length
	finish := page*length + length - 1
	rows, err := c.db.Query(SqlString + " WHERE id BETWEEN " +
		strconv.Itoa(start) + " AND " + strconv.Itoa(finish) + ";")
	if err != nil {
		return
	}
	defer rows.Close()

	// Iterate through results and create list of businesses
	for rows.Next() {
		var b Business
		err = rows.Scan(&b.Id, &b.Uuid, &b.Name, &b.Address, &b.Address2, &b.City, &b.Zip, &b.Country, &b.Phone, &b.Website, &b.Created_at)
		if err != nil {
			return
		}
		err = rows.Err()
		if err != nil {
			return
		}
		bs = append(bs, b)
	}

	// Catch if there are no results.
	if len(bs) == 0 {
		err = errors.New("No results found!")
	}
	return
}

// BusinessesHandler is the business end of the businesses handler.
func (c *Context) BusinessesHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the single id= case
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) > 0 {
		i, _ := strconv.Atoi(id)
		b, err := c.GetBusiness(i)
		if (err != nil) || (b == Business{}) {
			http.Error(w, err.Error(), 404)
		}
		o := FormatJson(b)
		o.WriteTo(w)
		return
	}

	// Parse query params
	length := 50
	lengthParam := vars["length"]
	if len(lengthParam) > 0 {
		length, _ = strconv.Atoi(lengthParam)
	}
	page := 0
	pageParam := vars["page"]
	if len(pageParam) > 0 {
		page, _ = strconv.Atoi(pageParam)
	}

	// Get businesses based on query and output
	bs, err := c.GetBusinesses(page, length)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	m := Meta{Page: page, Length: length}
	all := All{Businesses: bs, Meta: m}
	o := FormatJson(all)
	o.WriteTo(w)
}

// SetupDatabase loads the sqlite database from the specified file
func (c *Context) SetupDatabase() {
	var err error
	c.db, err = sql.Open("sqlite3", "businesses.sql")
	if err != nil {
		panic(err)
	}
}

// SetupRoutes sets up the specific routes for the API
func (c *Context) SetupRoutes() {
	c.router = mux.NewRouter()
	c.router.HandleFunc("/businesses", c.BusinessesHandler).Methods("GET").Queries("length", "{length:[0-9]+}", "page", "{page:[0-9]+}")
	c.router.HandleFunc("/businesses", c.BusinessesHandler).Methods("GET").Queries("page", "{page:[0-9]+}")
	c.router.HandleFunc("/businesses", c.BusinessesHandler).Methods("GET").Queries("length", "{length:[0-9]+}")
	c.router.HandleFunc("/businesses", c.BusinessesHandler).Methods("GET")
	c.router.HandleFunc("/businesses/", c.BusinessesHandler).Methods("GET")
	c.router.HandleFunc("/businesses/{id:[0-9]+}", c.BusinessesHandler).Methods("GET")
}

func main() {
	c := Context{}
	c.SetupDatabase()
	c.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", c.router))
}
