package handlers

import (
	"github.com/EtienneBerube/cat-scribers/internal/services"
	"log"
)

func HandleMonthlyPayments() {
	allUsers, err := services.GetAllUsers()
	if err != nil {
		log.Printf("CRON: Couldnt get all users. Error: %s", err.Error())
	}
	for _, user := range allUsers {
		services.PaySubscription(&user)
	}
}
