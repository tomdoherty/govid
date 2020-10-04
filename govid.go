package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"io"
	"io/ioutil"

	"encoding/json"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type covidRecord struct {
	DateRep                                          string `json:"dateRep"`
	Day                                              string `json:"day"`
	Month                                            string `json:"month"`
	Year                                             string `json:"year"`
	Cases                                            int    `json:"cases"`
	Deaths                                           int    `json:"deaths"`
	CountriesAndTerritories                          string `json:"countriesAndTerritories"`
	GeoID                                            string `json:"geoId"`
	CountryterritoryCode                             string `json:"countryterritoryCode"`
	PopData2019                                      int    `json:"popData2019"`
	ContinentExp                                     string `json:"continentExp"`
	CumulativeNumberFor14DaysOfCOVID19CasesPer100000 string `json:"Cumulative_number_for_14_days_of_COVID-19_cases_per_100000"`
}

type covidRecords struct {
	Records []struct {
		covidRecord
	} `json:"records"`
}

func loadInput(file string, f func(r struct{ covidRecord })) error {
	infile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("loadInput(%s): %e", file, err)
	}
	defer infile.Close()

	return parseInput(infile, f)
}

func parseInput(handle io.Reader, f func(r struct{ covidRecord })) error {
	var records covidRecords

	s, err := ioutil.ReadAll(handle)

	if err != nil {
		return fmt.Errorf("parseInput(%v): %s", handle, err)
	}

	if err := json.Unmarshal(s, &records); err != nil {
		return fmt.Errorf("parseInput(%v) json.Unmarshal: %s", handle, err)
	}

	for _, record := range records.Records {
		f(record)
	}
	return nil
}

func main() {
	// docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -p 3306:3306 -d mysql:latest
	// brew install mysql-client
	// /usr/local/opt/mysql-client/bin/mysql -h127.0.0.1 -uroot -pmy-secret-pw
	// create database covid;
	// create table covid (DateRep date, Cases int, Deaths int, GeoID varchar(100), Population int);

	db, err := sql.Open("mysql", "root:my-secret-pw@/covid")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO covid VALUES( ?, ?, ?, ?, ? )")

	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmtIns.Close()

	loadInput("covid.json", func(r struct{ covidRecord }) {
		time, err := time.Parse("02/01/2006", r.DateRep)
		_, err = stmtIns.Exec(time, r.Cases, r.Deaths, r.GeoID, r.PopData2019)
		if err != nil {
			log.Fatal(err)
		}
	})
}
