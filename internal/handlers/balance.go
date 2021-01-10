package handlers

import (
	"github.com/EtienneBerube/cat-scribers/internal/services"
	"log"
)

// HandleMonthlyPayments handles the payments of all users registered. (It is recommended to call this function in a CRON job)
func HandleMonthlyPayments() {
	allUsers, err := services.GetAllUsers()
	if err != nil {
		log.Printf("CRON: Couldnt get all users. Error: %s", err.Error())
	}
	for _, user := range allUsers {
		services.PaySubscription(&user)
	}
}
