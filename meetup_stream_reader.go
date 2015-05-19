package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Event struct {
	Description     string                 `json:"description"`
	Duration        int                    `json:"duration"`
	EventUrl        string                 `json:"event_url"`
	Fee             map[string]interface{} `json:"fee"`
	Group           map[string]interface{} `json:"group"`
	Id              string                 `json:"id"`
	MTime           int                    `json:"mtime"`
	Name            string                 `json:"name"`
	PaymentRequired string                 `json:"payment_required"`
	PhotoUrl        string                 `json:"photo_url"`
	RsvpLimit       int                    `json:"rsvp_limit"`
	Status          string                 `json:"status"`
	Time            int                    `json:"time"`
	UtcOffset       int                    `json:"utc_offset"`
	Venue           map[string]interface{} `json:"venue"`
	VenueVisibility string                 `json:"venue_visibility"`
	YesRsvpCount    int                    `json:"yes_rsvp_count"`
}

func (e Event) String() string {
	return fmt.Sprintf("%s - %s :: %s\n", e.Id, e.Name, e.Status)
}

func readEvents(events chan []byte) {
	for {
		var event Event
		e := <-events
		if err := json.Unmarshal(e, &event); err != nil {
			log.Println(err)
		}
		if event.Status == "upcoming" {
			fmt.Print(event)
		}
	}
}

func readStream() chan []byte {
	events := make(chan []byte, 25)
	go func() {
		for {
			resp, err := http.Get("http://stream.meetup.com/2/open_events?since_count=10")
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
