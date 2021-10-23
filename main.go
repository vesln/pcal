package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

type Category struct {
	Name      string
	Matcher   string
	Duration  float64
	Attendees int
}

func (c *Category) isMatching(e *gocal.Event) bool {
	if len(c.Matcher) > 0 && !strings.Contains(e.Summary, c.Matcher) {
		return false
	}

	if c.Attendees != len(e.Attendees) {
		return false
	}

	return true
}

func main() {
	var categories = []Category{
		Category{
			Name:      "Information Processing",
			Matcher:   "Information Processing",
			Attendees: 0,
		},
		Category{
			Name:      "Engineering",
			Matcher:   "Engineering",
			Attendees: 0,
		},
		Category{
			Name:      "Ops",
			Matcher:   "Ops",
			Attendees: 0,
		},
		Category{
			Name:      "Strategic",
			Matcher:   "Strategic",
			Attendees: 0,
		},
		Category{
			Name:      "Learning",
			Matcher:   "Learning",
			Attendees: 0,
		},
		Category{
			Name:      "Team", // should also match events with 2 or more attendees
			Matcher:   "Team",
			Attendees: 0,
		},
		Category{
			Name:      "1:1s", // Management
			Attendees: 2,
		},
		Category{
			Name:      "Recruiting", // recruiting@angel.co
			Matcher:   "Recruiting",
			Attendees: 0,
		},
	}

	// TODO: make it an input
	f, _ := os.Open("/Users/vesln/Downloads/ves@angel.co.ical/ves@angel.co.ics")
	defer f.Close()

	// TODO: make it an input
	start, end := time.Now().Add(-7*24*time.Hour), time.Now()

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	c.Parse()

	durations := make(map[string]float64)

	for _, e := range c.Events {
		for i, _ := range categories {
			if categories[i].isMatching(&e) {
				durations[categories[i].Name] += e.End.Sub(*e.Start).Minutes()
				break
			}
		}
	}

	for name, duration := range durations {
		fmt.Println(name, duration, "minutes")
	}
}
