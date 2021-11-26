package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Category struct {
	Name      string `json:"name"`
	Id        int    `json:"id"`
	Shortname string `json:"shortname"`
}

type Group struct {
	JoinMode string   `json:"join_mode"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Name     string   `json:"name"`
	GroupLon float32  `json:"group_lon"`
	GroupLat float32  `json:"group_lat"`
	Id       int      `json:"id"`
	URLName  string   `json:"urlname"`
	Category Category `json:"category"`
}

type Event struct {
	Description     string                 `json:"description"`
	Duration        int                    `json:"duration"`
	EventUrl        string                 `json:"event_url"`
	Fee             map[string]interface{} `json:"fee"`
	Group           Group                  `json:"group"`
	Id              string                 `json:"id"`
	MTime           int64                  `json:"mtime"`
	Name            string                 `json:"name"`
	PaymentRequired string                 `json:"payment_required"`
	PhotoUrl        string                 `json:"photo_url"`
	RsvpLimit       int                    `json:"rsvp_limit"`
	Status          string                 `json:"status"`
	Time            int64                  `json:"time"`
	UtcOffset       int                    `json:"utc_offset"`
	Venue           map[string]interface{} `json:"venue"`
	VenueVisibility string                 `json:"venue_visibility"`
	YesRsvpCount    int                    `json:"yes_rsvp_count"`
}

const (
	url        = "http://stream.meetup.com/2/open_events?since_count=10"
	dataFormat = "2006-01-02 15:04"
	status     = "upcoming"
)

func (e *Event) format() string {
	return fmt.Sprintf("%s - %s @ %s %s - %s",
		e.Group.Name,
		strings.TrimSpace(e.Name),
		strings.ToUpper(e.Group.Country),
		e.Group.City,
		time.Unix(0, e.Time*1000000).Format(dataFormat))
}

func readEvents(events chan []byte) {
	for {
		var event Event
		e := <-events
		if err := json.Unmarshal(e, &event); err != nil {
			log.Println(err)
		}
		if event.Status == status {
			log.Println(event.format())
		}
	}
}

func readStream() chan []byte {
	events := make(chan []byte, 25)
	go func() {
		for {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				break
			}
			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					fmt.Println(err)
					break
				}
				events <- line
			}
		}
	}()
	return events
}

func main() {
	fmt.Println("Loading Meetup Event Stream")

	var wg sync.WaitGroup
	wg.Add(1)

	events := readStream()
	go readEvents(events)

	wg.Wait()
}
