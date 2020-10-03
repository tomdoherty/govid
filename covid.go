package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"io/ioutil"

	"encoding/json"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type CovidRecord struct {
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
type CovidRecords struct {
	Records []struct {
		CovidRecord
	} `json:"records"`
}

func main() {
	f, err := os.Open("covid.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	// docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -p 3306:3306 -d mysql:latest
	// brew install mysql-client
	// /usr/local/opt/mysql-client/bin/mysql -h127.0.0.1 -uroot -pmy-secret-pw
	// create table covid (DateRep date, Cases int, Deaths int, GeoID varchar(100), Population int);
	db, err := sql.Open("mysql", "root:my-secret-pw@/covid")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO covid VALUES( ?, ?, ?, ?, ? )")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	var records CovidRecords
	if err := json.Unmarshal(s, &records); err != nil {
		log.Fatal(err)
	}
	totals := map[int]struct{ CovidRecord }{}
	var highest struct{ CovidRecord }

	for _, record := range records.Records {
		time, err := time.Parse("02/01/2006", record.DateRep)
		if err != nil {
			log.Fatal(err)
		}
		if record.Deaths > highest.Deaths {
			highest = record
		}
		totals[record.Deaths] = record
		//		log.Println(record.DateRep, record.Deaths, time)
		//DateRep varchar(10), Cases int, Deaths int, GeoID varchar(4), Population int )
		_, err = stmtIns.Exec(time, record.Cases, record.Deaths, record.GeoID, record.PopData2019)
		if err != nil {
			log.Fatal(err)
		}
	}

	keys := make([]int, 0, len(totals))

	for _, i := range totals {
		keys = append(keys, i.Deaths)
	}
	sort.Ints(keys)
	for _, i := range keys {
		fmt.Printf("%s\t%s\t%d\t%d\n", totals[i].DateRep, totals[i].GeoID, totals[i].PopData2019, totals[i].Deaths)
	}

	//	log.Printf("%+v", highest)

}
