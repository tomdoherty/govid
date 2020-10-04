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

	"github.com/tomdoherty/govid/pkg/govid"
)

func loadInput(file string, f func(r struct{ govid.CovidRecord })) error {
	infile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("loadInput(%s): %e", file, err)
	}
	defer infile.Close()

	return parseInput(infile, f)
}

func parseInput(handle io.Reader, f func(r struct{ govid.CovidRecord })) error {
	var records govid.CovidRecords

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

	loadInput("covid.json", func(r struct{ govid.CovidRecord }) {
		time, err := time.Parse("02/01/2006", r.DateRep)
		fmt.Println(time, r.Cases, r.Deaths, r.GeoID, r.PopData2019)

		_, err = stmtIns.Exec(time, r.Cases, r.Deaths, r.GeoID, r.PopData2019)
		if err != nil {
			log.Fatal(err)
		}
	})
}
