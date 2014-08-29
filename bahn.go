package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

func main() {

	if len(os.Args) < 3 {
		usage()
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', 0)

	var from, to string = os.Args[1], os.Args[2]

	var Url *url.URL
	Url, err := url.Parse("http://reiseauskunft.bahn.de/bin/query.exe/dn")

	parameters := url.Values{}
	parameters.Add("S", from)
	parameters.Add("Z", to)
	parameters.Add("start", "1")
	Url.RawQuery = parameters.Encode()

	doc, err := goquery.NewDocument(Url.String())
	if err != nil {
		log.Fatal(err)
	}

	firstrows := []string{}
	lastrows := []string{}
	result := doc.Find(".result")

	result.Find(".firstrow").Each(func(i int, s *goquery.Selection) {
		time := strings.TrimSpace(strings.Trim(s.Find("td.time").Text(), "\n"))[0:5]
		date := strings.TrimSpace(strings.Trim(s.Find("td.date").Text(), "\n"))
		changes := strings.TrimSpace(strings.Trim(s.Find("td.changes").Text(), "\n"))
		duration := strings.TrimSpace(strings.Trim(s.Find("td.duration").Text(), "\n"))
		product := strings.TrimSpace(strings.Trim(s.Find("td.products").Text(), "\n"))
		station := strings.TrimSpace(strings.Trim(s.Find("td.station > div.resultDep").Text(), "\n"))
		firstrows = append(firstrows, station+"\t"+date+"\t"+time+"\t"+duration+"\t"+changes+"\t"+product)
	})

	result.Find(".last").Each(func(i int, s *goquery.Selection) {
		time := strings.TrimSpace(strings.Trim(s.Find("td.time").Text(), "\n"))
		date := strings.TrimSpace(strings.Trim(s.Find("td.date").Text(), "\n"))
		station := strings.TrimSpace(strings.Trim(s.Find("td.station.stationDest").Text(), "\n"))
		lastrows = append(lastrows, station+"\t"+date+"\t"+time+"\t\t\t")
	})

	fmt.Fprintln(w, "Station\tDate\tTime\tDuration\tChanges\tType")
	for i := range firstrows {
		fmt.Fprintln(w, firstrows[i])
		fmt.Fprintln(w, lastrows[i])
		fmt.Fprintln(w, "\t\t\t\t\t\t")
	}

	fmt.Fprintln(w)
}

func usage() {
	fmt.Printf("usage: %s [from] [to]\n", os.Args[0])
	os.Exit(2)
}
