package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tomdoherty/govid/pkg/govid"
)

var (
	LoadInput = loadInput
)

func TestLoadInput(t *testing.T) {
	t.Run("Test loadInput", func(t *testing.T) {
		var output bytes.Buffer
		if err := parseInput(strings.NewReader(`{
   "records" : [
      {
         "dateRep" : "02/10/2020",
         "day" : "02",
         "month" : "10",
         "year" : "2020",
         "cases" : 17,
         "deaths" : 0,
         "countriesAndTerritories" : "Afghanistan",
         "geoId" : "AF",
         "countryterritoryCode" : "AFG",
         "popData2019" : 38041757,
         "continentExp" : "Asia",
         "Cumulative_number_for_14_days_of_COVID-19_cases_per_100000" : "1.08564912"
      }
   ]
}
`), func(r struct{ govid.CovidRecord }) {
			time, err := time.Parse("02/01/2006", r.DateRep)
			if err != nil {
				t.Fatal("Failed to parse time")
			}
			fmt.Fprintf(&output, "%s,%d,%d,%s,%d", time, r.Cases, r.Deaths, r.GeoID, r.PopData2019)

		}); err != nil {
			t.Errorf("loadInput() error = %v", err)
		}

		if strings.Contains("geoID", output.String()) {
			t.Fatal("Failed to parse input")
		}
	})

}
