package main

import (
	"log"
	"os"

	"github.com/tomdoherty/govid"
)

func main() {
	var output govid.SQLWriter
	//var output govid.LogWriter
	//if err := output.Init(os.Stdout); err != nil {
	if err := output.Init(); err != nil {
		log.Fatal(err)
	}
	govid.Filter(os.Stdin, &output)
}
