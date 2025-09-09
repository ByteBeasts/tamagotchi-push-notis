package utils

import (
	"encoding/json"
	"log"
	"strings"
	"tamagotchi-push-notis/common"
)

// extractUsernameFromEmail extracts the username part from an email address
func extractUsernameFromEmail(email string) string {
	trimmed := strings.TrimSpace(email)

	// Remove leading single quote if present
	trimmed = strings.TrimPrefix(trimmed, "'")

	// Remove spaces
	cleaned := strings.ReplaceAll(trimmed, " ", "")

	if strings.Contains(cleaned, "@") {
		parts := strings.Split(cleaned, "@")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// If no @ symbol, return the cleaned address
	return cleaned
}

// addWalletPrefix adds "0x" prefix to an address if it doesn't already have it
func addWalletPrefix(address string) string {
	if strings.HasPrefix(address, "0x") {
		return address
	}
	return "0x" + address
}

// isValidWalletAddress checks if an address is a valid wallet address (42 characters including 0x)
func isValidWalletAddress(address string) bool {
	return len(address) == 42
}

// BatchedAddresses batches the addresses into smaller arrays of the given size
func BatchedAddresses(addresses []string, batchSize int) [][]string {
	var batchedAddresses [][]string
	var currentBatch []string

	for _, address := range addresses {
		// Only add non-empty addresses
		if strings.TrimSpace(address) != "" {
			currentBatch = append(currentBatch, address)

			// When batch is full, add it to the result and start a new batch
			if len(currentBatch) == batchSize {
				batchedAddresses = append(batchedAddresses, currentBatch)
				currentBatch = []string{}
			}
		}
	}

	// Add the remaining addresses if any
	if len(currentBatch) > 0 {
		batchedAddresses = append(batchedAddresses, currentBatch)
	}

	return batchedAddresses
}

// CleanAddresses cleans the addresses by removing the part after the @ symbol,
// adding 0x prefix, and filtering out invalid wallet addresses
func CleanAddresses(addresses []string) []string {
	var cleanedAddresses []string

	for _, address := range addresses {
		// Extract username from email
		username := extractUsernameFromEmail(address)

		// Add wallet prefix
		walletAddress := addWalletPrefix(username)

		// Only keep valid wallet addresses (42 characters)
		if isValidWalletAddress(walletAddress) {
			cleanedAddresses = append(cleanedAddresses, walletAddress)
		}
	}

	return cleanedAddresses
}

// PreparePayload prepares the payload for the API
func PreparePayload(appID string, addresses []string, title string, message string, miniAppPath string) []byte {
	payload := common.Payload{
		AppID:       appID,
		Addresses:   addresses,
		Title:       title,
		Message:     message,
		MiniAppPath: miniAppPath,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}
	return data
}
