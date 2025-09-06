package main

import (
	"fmt"
	"os"
	"tamagotchi-push-notis/http"
	"tamagotchi-push-notis/parser"
	"tamagotchi-push-notis/utils"
)

// var (
// 	messages = map[string]string{
// 		"Play!":  "Be the top 1, keep playing minigames!",
// 		"Care":   "Remember to clean your Beasts!",
// 		"Feed":   "Don't forget to feed your Beasts!",
// 		"Hungry": "Your Beast might be hungry!",
// 		"Happy":  "Keep it up one day is one year!",
// 	}
// )

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

	// Check if we have at least one column
	if len(csvData.Headers) < 1 {
		fmt.Printf("Error: CSV must have at least 1 column, but found %d\n", len(csvData.Headers))
		return
	}

	// Clean the addresses from the first (and only) column
	// cleanedAddresses := utils.CleanAddresses(csvData.GetColumn(csvData.Headers[0]))

	// Use dummy addresses for testing
	dummyAddresses := []string{
		"0x7b9e2692aa4b72e325808611d10ca128b8bc8eb6",
	}

	title := "Play!"
	message := "Be the top 1, keep playing minigames!"

	appID := os.Getenv("APP_ID")
	if appID == "" {
		fmt.Println("Error: APP_ID environment variable is required")
		return
	}

	// Prepare the payload with dummy addresses for testing
	payload := utils.PreparePayload(appID, dummyAddresses, title, message, "worldapp://mini-app?app_id="+appID)
	fmt.Println("Payload:", string(payload))

	// Get URL and bearer token from environment variables
	worldUrl := os.Getenv("WORLD_URL")
	worldBearer := os.Getenv("WORLD_BEARER")

	if worldUrl == "" || worldBearer == "" {
		fmt.Println("Error: WORLD_URL and WORLD_BEARER environment variables are required")
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

	// Print detailed response information
	fmt.Println("=== HTTP POST Response Details ===")
	fmt.Printf("Status Code: %d\n", response.StatusCode)
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Content Length: %d\n", response.ContentLength)
	fmt.Printf("Content Type: %s\n", response.Header.Get("Content-Type"))
	fmt.Printf("Response Headers:\n")
	for key, values := range response.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Read and print response body
	bodyBytes := make([]byte, 1024) // Read first 1024 bytes
	n, err := response.Body.Read(bodyBytes)
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("Error reading response body: %v\n", err)
	} else {
		fmt.Printf("Response Body (first %d bytes): %s\n", n, string(bodyBytes[:n]))
	}

	fmt.Println("=== End Response Details ===")
}
