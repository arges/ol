// (c) 2016 Written by Chris J Arges <christopherarges@gmail.com>

package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// url/status should match
var url_tests = []string{"/businesses/", "/businesses?page=1", "/businesses/1",
	"/businesses?length=1", "/businesses?page=0&length=1", "/", "/blah",
	"/businesses/-1", "/businesses/23432432432432"}
var status_tests = []int{http.StatusOK, http.StatusOK, http.StatusOK,
	http.StatusOK, http.StatusOK, http.StatusNotFound, http.StatusNotFound,
	http.StatusNotFound, http.StatusNotFound}

// TestBusinessHandler tests various url patterns for status codes
func TestBusinessesHandler(t *testing.T) {
	c := Context{}
	c.SetupDatabase()
	c.SetupRoutes()
	server := httptest.NewServer(c.router)
	defer server.Close()
	for index, url := range url_tests {
		resp, _ := http.Get(server.URL + url)
		if resp.StatusCode != status_tests[index] {
			t.Error("Invalid response code for: " + url + " got: " +
				resp.Status + " wanted: " + strconv.Itoa(status_tests[index]))
		}
	}
}
