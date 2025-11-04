package main

import (
	"library_management/controllers"
	"library_management/services"
	"library_management/models"
	"fmt"
)

func main() {
	fmt.Println("===== Welcome to MicrQ Library ===== \n")

	library := services.Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}

	controllers.Console(&library)
}
