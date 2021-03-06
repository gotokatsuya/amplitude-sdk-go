package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gotokatsuya/amplitude-sdk-go/amplitude"
)

var (
	apiKey = os.Getenv("API_KEY")
)

func main() {
	ctx := context.Background()

	cli, err := amplitude.NewClient(apiKey, http.DefaultClient)
	if err != nil {
		panic(err)
	}

	// log event
	eventRes, httpRes, err := cli.LogEvent(ctx, &amplitude.LogEventRequest{
		Events: []amplitude.Event{
			{
				UserID:    "john_doe@gmail.com",
				EventType: "watch_tutorial",
				UserProperties: map[string]interface{}{
					"Cohort": "Test A",
				},
				Country: "United States",
				IP:      "127.0.0.1",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	switch httpRes.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		log.Printf("Event: %v\n", eventRes)
	default:
		log.Printf("Error: %s\n", eventRes.Error)
		return
	}

	// identify
	httpRes, err = cli.Identify(ctx, &amplitude.IdentifyRequest{
		Identifications: []amplitude.Identification{
			{
				UserID: "john_doe@gmail.com",
				UserProperties: map[string]interface{}{
					"Cohort": "Test B",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	switch httpRes.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		log.Printf("Identify Success\n")
	default:
		log.Printf("Identify Error\n")
		return
	}
}
