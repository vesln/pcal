package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/apognu/gocal"
	log "github.com/sirupsen/logrus"
)

const (
	Accepted = "ACCEPTED"
)

type Event struct {
	Summary   string
	Category  string
	Duration  time.Duration
	StartDate time.Time
}

type Category struct {
	Name          string `json:"name"`
	Prefix        string `json:"prefix"`
	MinAttendees  int    `json:"min"`
	MaxAttendees  int    `json:"max"`
	AttendeeEmail string `json:"email"`
	Ignore        bool   `json:"ignore"`
}

func (c Category) isMatching(e gocal.Event) bool {
	// Check if we have a prefix that's matching
	if len(c.Prefix) > 0 && !strings.HasPrefix(e.Summary, c.Prefix) {
		return false
	}

	totalAttendees := len(e.Attendees)

	// Check if the attendees count is within the bounds
	if totalAttendees < c.MinAttendees || totalAttendees > c.MaxAttendees {
		return false
	}

	// Ignore all events that I didn't accept
	if totalAttendees > 0 {
		for _, a := range e.Attendees {
			if a.Cn == *email && a.Status != Accepted {
				return false
			}
		}
	}

	// Check for special emails
	if len(c.AttendeeEmail) > 0 {
		ok := false

		for _, a := range e.Attendees {
			if a.Cn == c.AttendeeEmail {
				ok = true
			}
		}

		if !ok {
			return false
		}
	}

	return true
}

func parseIcs(icsPath string, startDate string, endDate string) []gocal.Event {
	// Format the start and end dates
	start, err := time.Parse(time.RFC822, startDate)
	if err != nil {
		log.Fatal(err)
	}

	end, err := time.Parse(time.RFC822, endDate)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the ics file
	icsFile, _ := os.Open(icsPath)
	defer icsFile.Close()

	c := gocal.NewParser(icsFile)
	c.Start = &start
	c.End = &end
	c.Parse()

	return c.Events
}

func parseCategories() []Category {
	usr, _ := user.Current()

	jsonFile, err := os.Open(fmt.Sprintf("%s/.config/pcal/categories.json", usr.HomeDir))
	defer jsonFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	categories := []Category{}

	json.Unmarshal(byteValue, &categories)

	return categories
}

func makeEvents(events []gocal.Event, categories []Category) []Event {
	var ret []Event

	for _, e := range events {
		found := false

		for _, c := range categories {
			if !c.isMatching(e) {
				continue
			}

			log.Debug("[", c.Name, "]", " ", e.Summary)
			found = true

			if !c.Ignore {
				duration := e.End.Sub(*e.Start)
				ret = append(ret, Event{Summary: e.Summary, Duration: duration, Category: c.Name, StartDate: *e.Start})
			}

			break
		}

		if !found {
			log.Debug("[IGNORED] ", e.Summary)
		}
	}

	return ret
}
