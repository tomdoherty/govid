package govid

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"database/sql"

	// optional mysql
	_ "github.com/go-sql-driver/mysql"
)

// CovidRecord is an individual record
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

// CovidRecords is an array of CovidRecord
type CovidRecords struct {
	Records []CovidRecord `json:"records"`
}

// RecordWriter outputs our records in normalised form
type RecordWriter interface {
	WriteRecord(r *CovidRecord) error
}

// SQLWriter is a RecordWriter which writes to MySQL
type SQLWriter struct {
	db   *sql.DB
	stmt *sql.Stmt
}

// WriteRecord to MySQL
func (w *SQLWriter) WriteRecord(r *CovidRecord) error {
	time, err := time.Parse("02/01/2006", r.DateRep)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.stmt.Exec(time, r.Cases, r.Deaths, r.GeoID, r.PopData2019)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Init database connection/prepared query
func (w *SQLWriter) Init() (err error) {
	w.db, err = sql.Open("mysql", "root:my-secret-pw@/covid")
	if err != nil {
		return err
	}

	w.stmt, err = w.db.Prepare("INSERT INTO covid VALUES( ?, ?, ?, ?, ? )")
	if err != nil {
		return err
	}

	return nil
}

// LogWriter is a RecordWriter which writes to an io.Writer
type LogWriter struct {
	w io.Writer
}

// WriteRecord to io.Writer
func (w *LogWriter) WriteRecord(r *CovidRecord) error {
	time, err := time.Parse("02/01/2006", r.DateRep)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintf(w.w, "%s\t%d\t%d\t%s%d\n", time.String(), r.Cases, r.Deaths, r.GeoID, r.PopData2019)
	return err
}

// Init LogWriter
func (w *LogWriter) Init(r io.Writer) (err error) {
	w.w = r
	return nil
}

// Filter takes an io.Reader and writes CovidRecord to a RecordWriter
func Filter(r io.Reader, w RecordWriter) {
	var v CovidRecords
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		log.Fatal(err)
	}
	for _, record := range v.Records {
		w.WriteRecord(&record)
	}
}
