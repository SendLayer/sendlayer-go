package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sendlayer/sendlayer-go/sendlayer"
)

func main() {
	apiKey := os.Getenv("SENDLAYER_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set SENDLAYER_API_KEY environment variable")
		os.Exit(1)
	}

	sl := sendlayer.New(apiKey)

	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()
	event := "accepted"
	retrieveCount := 10
	resp, err := sl.Events.Get(&startDate, &endDate, event, "", nil, &retrieveCount)
	if err != nil {
		fmt.Printf("Error getting events: %v\n", err)
		return
	}

	fmt.Printf("📊 Found %d events in the last 7 days:\n", resp.TotalRecords)

	if len(resp.Events) == 0 {
		fmt.Println("No events found in the specified time range.")
		return
	}
	fmt.Printf("Events: %v\n", resp)
}
