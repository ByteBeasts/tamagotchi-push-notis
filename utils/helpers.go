package utils

import (
	"log"
	"math/rand"
	"os"
	"tamagotchi-push-notis/http"
	"tamagotchi-push-notis/parser"
	"time"
)

const BatchSize = 500

var (
	Messages = map[string]string{
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

	WorldAppPath = "worldapp://mini-app?app_id="
)

// RandomMessage returns a random title and message from the messages map
func RandomMessage() (string, string) {
	keys := make([]string, 0, len(Messages))
	for key := range Messages {
		keys = append(keys, key)
	}
	randomKey := keys[rand.Intn(len(keys))]
	return randomKey, Messages[randomKey]
}

// ValidateEnvironmentVariables checks if all required environment variables are set
func ValidateEnvironmentVariables() (string, string, string, string, string) {
	cavosUrl := os.Getenv("CAVOS_URL")
	cavosBearer := os.Getenv("CAVOS_BEARER")
	worldUrl := os.Getenv("WORLD_URL")
	worldBearer := os.Getenv("WORLD_BEARER")
	appID := os.Getenv("APP_ID")

	if cavosUrl == "" || cavosBearer == "" {
		log.Fatalf("Error: CAVOS_URL and CAVOS_BEARER environment variables are required")
	}

	if worldUrl == "" || worldBearer == "" {
		log.Fatalf("Error: WORLD_URL and WORLD_BEARER environment variables are required")
	}

	if appID == "" {
		log.Fatalf("Error: APP_ID environment variable is required")
	}

	return cavosUrl, cavosBearer, worldUrl, worldBearer, appID
}

// FetchCSVData makes HTTP request to Cavos API and parses CSV response
func FetchCSVData(cavosUrl, cavosBearer string) *parser.CSVData {
	// Make HTTP request
	requester := http.NewRequester(cavosUrl, cavosBearer)
	response, err := requester.Request()
	if err != nil {
		log.Fatalf("Error making request: %v\n", err)
	}
	defer response.Body.Close()

	// Log response status for debugging
	log.Printf("CAVOS HTTP Response Status: %s (Code: %d)\n", response.Status, response.StatusCode)

	// Parse CSV response
	csvParser := parser.NewParser(response.Body)
	csvData, err := csvParser.Parse()
	if err != nil {
		log.Fatalf("Error parsing CSV: %v\n", err)
	}

	// Check if we have at least one column
	if len(csvData.Headers) < 1 {
		log.Fatalf("Error: CSV must have at least 1 column, but found %d\n", len(csvData.Headers))
	}

	return csvData
}

// ProcessAddresses extracts and cleans addresses from CSV data
func ProcessAddresses(csvData *parser.CSVData) []string {
	// Get raw addresses from CSV
	rawAddresses := csvData.GetColumn(csvData.Headers[0])

	// Clean the addresses from the first (and only) column
	cleanedAddresses := CleanAddresses(rawAddresses)

	return cleanedAddresses
}

// sendBatchNotification sends a notification batch to the World API
func sendBatchNotification(batch []string, appID, worldUrl, worldBearer string) {
	log.Printf("Batch size: %d\n", len(batch))

	// Pick a random message from the map
	title, message := RandomMessage()

	// Prepare the payload
	payload := PreparePayload(appID, batch, title, message, WorldAppPath+appID)

	// Post the payload
	poster := http.NewPoster(worldUrl, worldBearer, payload)

	// Check if the payload is posted
	response, err := poster.Post()
	if err != nil {
		log.Fatalf("Error posting: %v\n", err)
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

// SendAllNotifications processes all address batches and sends notifications
func SendAllNotifications(addresses []string, appID, worldUrl, worldBearer string) {
	// Batch the addresses
	batchedAddresses := BatchedAddresses(addresses, BatchSize)

	// Send notifications for each batch
	for _, batch := range batchedAddresses {
		sendBatchNotification(batch, appID, worldUrl, worldBearer)
	}
}

// RunNotificationProcess executes the complete notification process
func RunNotificationProcess() {
	log.Println("Starting notification process...")

	// Validate environment variables
	cavosUrl, cavosBearer, worldUrl, worldBearer, appID := ValidateEnvironmentVariables()

	// Fetch and parse CSV data from Cavos API
	csvData := FetchCSVData(cavosUrl, cavosBearer)

	// Process addresses from CSV
	cleanedAddresses := ProcessAddresses(csvData)

	// Send notifications to all batches using the WorldCoin API
	SendAllNotifications(cleanedAddresses, appID, worldUrl, worldBearer)

	log.Println("Notification process completed successfully!")
}

// StartScheduler starts a goroutine that runs the notification process every 9 hours
func StartScheduler() {
	log.Println("Starting notification scheduler (every 9 hours)...")

	// Run immediately on startup
	RunNotificationProcess()

	// Create a ticker for 9 hours
	ticker := time.NewTicker(9 * time.Hour)
	defer ticker.Stop()

	// Run in a goroutine to keep the main thread free
	go func() {
		for range ticker.C {
			log.Println("Scheduled notification process triggered...")
			RunNotificationProcess()
		}
	}()
}
