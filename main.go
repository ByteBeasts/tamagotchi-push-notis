package main

import (
	"fmt"
	"os"
	"tamagotchi-push-notis/http"
	"tamagotchi-push-notis/parser"
	"tamagotchi-push-notis/utils"
)

var (
	messages = []string{
		"Be the top 1, keep playing minigames!",
		"Remember to clean your Beasts!",
		"Don't forget to feed your Beasts!",
		"Your Beast might be hungry!",
		"Keep it up one day is one year!",
	}
)

func main() {
	// Get URL and bearer token from environment variables
	cavosUrl := os.Getenv("CAVOS_URL")
	cavosBearer := os.Getenv("CAVOS_BEARER")

	if cavosUrl == "" || cavosBearer == "" {
		fmt.Println("Error: CAVOS_URL and CAVOS_BEARER environment variables are required")
		return
	}

	// Make HTTP request
	requester := http.NewRequester(cavosUrl, cavosBearer)
	response, err := requester.Request()
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	// Parse CSV response
	csvParser := parser.NewParser(response.Body)
	csvData, err := csvParser.Parse()
	if err != nil {
		fmt.Println("Error parsing CSV:", err)
		return
	}

	// Clean the addresses
	cleanedAddresses := utils.CleanAddresses(csvData.GetColumn(csvData.Headers[1]))

	// Prepare the payload
	payload := utils.PreparePayload(csvData.Headers[0], cleanedAddresses, messages[0])


	// Get URL and bearer token from environment variables
	worldUrl := os.Getenv("WORLD_URL")
	worldBearer := os.Getenv("WORLD_BEARER")

	if worldUrl == "" || worldBearer == "" {
		fmt.Println("Error: CAVOS_URL and CAVOS_BEARER environment variables are required")
		return
	}

	// Post the payload
	poster := http.NewPoster(worldUrl, worldBearer, payload)

	// Check if the payload is posted
	response, err = poster.Post()
	if err != nil {
		fmt.Println("Error posting:", err)
		return
	}
	defer response.Body.Close()

	// Print the response
	fmt.Println("Response:", response)
}
