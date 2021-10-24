package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	email   *string
	icsPath *string
	start   *string
	end     *string
	debug   *bool
	format  *string
)

const (
	formatCsv   = "csv"
	formatAscii = "ascii"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func init() {
	icsPath = flag.String("path", "", "path to the .ics file")
	email = flag.String("email", "", "example@example.com")
	start = flag.String("start", "", "17 Oct 21 00:00 EET")
	end = flag.String("end", "", "17 Oct 21 00:00 EET")
	debug = flag.Bool("debug", false, "turn on debug output")
	format = flag.String("format", formatCsv, "format the output as csv or ascii")

	flag.Parse()

	log.SetOutput(os.Stdout)

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if !isFlagPassed("path") {
		log.Fatal("Please specify and .ics path")
	}

	if !isFlagPassed("email") {
		log.Fatal("Please specify an e-mail address")
	}

	if !isFlagPassed("start") {
		log.Fatal("Please specify a start date")
	}

	if !isFlagPassed("end") {
		log.Fatal("Please specify an end date")
	}
}

func main() {
	categories := parseCategories()
	icsEvents := parseIcs(*icsPath, *start, *end)

	events := makeEvents(icsEvents, categories)

	switch *format {
	case formatAscii:
		printAscii(events)
	case formatCsv:
		printCsv(events)
	default:
		log.Fatal("Unkonwn format")
	}
}
