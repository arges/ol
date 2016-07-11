all: business.sql
	go build

business.sql:
	./setup.sh

clean:
	rm -f businesses.sql engineering_project_businesses.csv ol
