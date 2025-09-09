package main

import (
	"log"
	"tamagotchi-push-notis/utils"
)

func main() {
	log.Println("ByteBeasts Tamagotchi Push Notifications Service Starting...")

	// Start the scheduler that runs every 9 hours
	utils.RunNotificationProcess()
}
