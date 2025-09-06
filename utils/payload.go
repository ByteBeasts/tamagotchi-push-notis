package utils

import (
	"encoding/json"
	"log"
	"strings"
	"tamagotchi-push-notis/common"
)

// CleanAddresses cleans the addresses by removing the part after the @ symbol
func CleanAddresses(addresses []string) []string {
	var cleanedAddresses []string
	for _, address := range addresses {
		// Trim whitespace first
		trimmed := strings.TrimSpace(address)

		// Extract the part before @ symbol
		if strings.Contains(trimmed, "@") {
			parts := strings.Split(trimmed, "@")
			if len(parts) > 0 {
				cleanedAddresses = append(cleanedAddresses, parts[0])
			}
		} else {
			// If no @ symbol, keep the original address
			cleanedAddresses = append(cleanedAddresses, trimmed)
		}
	}
	return cleanedAddresses
}

// PreparePayload prepares the payload for the API
func PreparePayload(appID string, addresses []string, title string, message string, miniAppPath string) []byte {
	payload := common.Payload{
		AppID:     appID,
		Addresses: addresses,
		Title:     title,
		Message:   message,
		MiniAppPath: miniAppPath,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}
	return data
}
