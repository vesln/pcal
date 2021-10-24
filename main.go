package main

import (
	"flag"
	"os"
	"time"

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

func GetStartDayOfWeek() time.Time {
	tm := time.Now()

	weekday := time.Duration(tm.Weekday())

	if weekday == 0 {
		weekday = 7
	}

	year, month, day := tm.Date()

	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}

func init() {
	icsPath = flag.String("path", "", "path to the .ics file")
	email = flag.String("email", "", "example@example.com")
	start = flag.String("start", GetStartDayOfWeek().Format(time.RFC822), "17 Oct 21 00:00 EET")
	end = flag.String("end", time.Now().Format(time.RFC822), "17 Oct 21 00:00 EET")
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
