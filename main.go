package main

import (
	"tamagotchi-push-notis/utils"
)

func main() {
	// Validate environment variables
	cavosUrl, cavosBearer, worldUrl, worldBearer, appID := utils.ValidateEnvironmentVariables()

	// Fetch and parse CSV data from Cavos API
	csvData := utils.FetchCSVData(cavosUrl, cavosBearer)

	// Process addresses from CSV
	cleanedAddresses := utils.ProcessAddresses(csvData)

	// Send notifications to all batches
	utils.SendAllNotifications(cleanedAddresses, appID, worldUrl, worldBearer)
}
