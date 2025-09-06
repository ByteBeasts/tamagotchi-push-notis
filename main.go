package main

import (
	"log"
	"math/rand"
	"os"
	"tamagotchi-push-notis/http"
	"tamagotchi-push-notis/parser"
	"tamagotchi-push-notis/utils"
)

const batchSize = 500

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

	worldAppPath = "worldapp://mini-app?app_id="
)

func randomMessage() (string, string) {
	keys := make([]string, 0, len(messages))
	for key := range messages {
		keys = append(keys, key)
	}
	return keys[rand.Intn(len(keys))], messages[keys[rand.Intn(len(keys))]]
}

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
	log.Printf("CAVOSHTTP Response Status: %s (Code: %d)\n", response.Status, response.StatusCode)

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

	// Clean the addresses from the first (and only) column
	cleanedAddresses := utils.CleanAddresses(rawAddresses)

	// Batch the addresses
	batchedAddresses := utils.BatchedAddresses(cleanedAddresses, batchSize)

	appID := os.Getenv("APP_ID")
	if appID == "" {
		log.Fatalf("Error: APP_ID environment variable is required")
		return
	}

	for _, batch := range batchedAddresses {
		log.Printf("Batch size: %d\n", len(batch))
		// Pick a random message from the map
		title, message := randomMessage()

		// Prepare the payload with dummy addresses for testing
		payload := utils.PreparePayload(appID, batch, title, message, worldAppPath+appID)

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
			log.Printf("WorldResponse Body (first %d bytes): %s\n", n, string(bodyBytes[:n]))
		}
	}
}
