package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	return fmt.Sprintf("%01dh %02dm", h, m)
}

func printAscii(events []Event) {
	var totalDuration time.Duration = 0
	var durations = make(map[string]time.Duration)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Total"})

	for _, event := range events {
		durations[event.Category] += event.Duration
		totalDuration += event.Duration
	}

	for name, d := range durations {
		table.Append([]string{name, formatDuration(d)})
	}

	table.SetFooter([]string{"", formatDuration(totalDuration)})

	table.Render()
}

func printCsv(events []Event) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	for _, event := range events {
		writer.Write([]string{event.StartDate.Format("2006-01-02"), event.Category, event.Summary, fmt.Sprintf("%.0f", event.Duration.Minutes())})
	}

	writer.Flush()
	fmt.Println(buffer.String())
}
