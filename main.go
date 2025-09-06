package main

import (
	"log"
	"math/rand"
	"os"
	"tamagotchi-push-notis/http"
	"tamagotchi-push-notis/parser"
	"tamagotchi-push-notis/utils"
)

var (
	messages = map[string]string{
		"Play!":      "🎮 Be the top 1 — keep playing minigames!",
		"Care":       "🧼 Remember to clean your Beasts!",
		"Feed":       "🍗 Don't forget to feed your Beasts!",
		"Hungry":     "🍽️ Your Beast might be hungry!",
		"Happy":      "😄 Keep it up — one day is one year!",
		"Sleep":      "🌙 Bedtime! Let your Beast recharge.",
		"Energy Low": "⚡️ Energy is low — a boost could help.",
		"Level Up":   "⭐️ So close to leveling up — one more game!",
		"Name Time":  "🏷️ Give your Beast a cool new name!",
		"Clean Up":   "🫧 Mud alert! Your Beast needs a bath.",
		"Miss You":   "👋 Your Beast misses you — come say hi!",
	}
)

func main() {
	// Get URL and bearer token from environment variables
	cavosUrl := os.Getenv("CAVOS_URL")
	cavosBearer := os.Getenv("CAVOS_BEARER")

	if cavosUrl == "" || cavosBearer == "" {
		log.Fatalf("Error: CAVOS_URL and CAVOS_BEARER environment variables are required")
		return
	}

	// Make HTTP request
	requester := http.NewRequester(cavosUrl, cavosBearer)
	response, err := requester.Request()
	if err != nil {
		log.Fatalf("Error making request: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Log response status for debugging
	log.Printf("HTTP Response Status: %s (Code: %d)\n", response.Status, response.StatusCode)

	// Parse CSV response
	csvParser := parser.NewParser(response.Body)
	csvData, err := csvParser.Parse()
	if err != nil {
		log.Fatalf("Error parsing CSV: %v\n", err)
		return
	}

	// Check if we have at least one column
	if len(csvData.Headers) < 1 {
		log.Fatalf("Error: CSV must have at least 1 column, but found %d\n", len(csvData.Headers))
		return
	}

	// Get raw addresses from CSV
	rawAddresses := csvData.GetColumn(csvData.Headers[0])
	// log.Printf("Raw addresses from CSV: %v\n", rawAddresses)

	// Clean the addresses from the first (and only) column
	cleanedAddresses := utils.CleanAddresses(rawAddresses)
	// log.Printf("Cleaned addresses: %v\n", cleanedAddresses)

	// Pick a random message from the map
	keys := make([]string, 0, len(messages))
	for key := range messages {
		keys = append(keys, key)
	}
	randomKey := keys[rand.Intn(len(keys))]
	title := randomKey
	message := messages[randomKey]

	appID := os.Getenv("APP_ID")
	if appID == "" {
		log.Fatalf("Error: APP_ID environment variable is required")
		return
	}

	// Prepare the payload with dummy addresses for testing
	payload := utils.PreparePayload(appID, cleanedAddresses, title, message, "worldapp://mini-app?app_id="+appID)
	log.Printf("Payload: %s\n", string(payload))

	// Get URL and bearer token from environment variables
	worldUrl := os.Getenv("WORLD_URL")
	worldBearer := os.Getenv("WORLD_BEARER")

	if worldUrl == "" || worldBearer == "" {
		log.Fatalf("Error: WORLD_URL and WORLD_BEARER environment variables are required")
		return
	}

	// Post the payload
	poster := http.NewPoster(worldUrl, worldBearer, payload)

	// Check if the payload is posted
	response, err = poster.Post()
	if err != nil {
		log.Fatalf("Error posting: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Read and print response body
	bodyBytes := make([]byte, 1024) // Read first 1024 bytes
	n, err := response.Body.Read(bodyBytes)
	if err != nil && err.Error() != "EOF" {
		log.Fatalf("Error reading response body: %v\n", err)
	} else {
		log.Printf("Response Body (first %d bytes): %s\n", n, string(bodyBytes[:n]))
	}
}
