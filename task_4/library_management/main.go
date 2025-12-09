package main

import (
	"fmt"
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	fmt.Println("===== Welcome to MicrQ Library =====")

	library := services.Library{
		Books:           make(map[int]models.Book),
		Members:         make(map[int]models.Member),
		ReservationChan: make(chan services.ReservationRequest, 10), // buffered channel
		CancelChans:     make(map[int]chan bool),
	}

	library.StartWorker()

	controllers.Console(&library)
}
