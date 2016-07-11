#!/bin/bash -x
set -e
sudo apt-get install sqlite3 golang-go
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
wget https://s3.amazonaws.com/ownlocal-engineering/engineering_project_businesses.csv.gz
gunzip engineering_project_businesses.csv.gz
echo '
create table businesses(id integer not null primary key,uuid text,name text,address text,address2 text,city text,state text,zip text,country text,phone text,website text,created_at text);
.mode csv
.import ./engineering_project_businesses.csv businesses
' | sqlite3 businesses.sql
