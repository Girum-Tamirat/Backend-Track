package main

import (
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
	"library_management/concurrency"
)

func main() {
	library := services.NewLibrary()
	library.AddBook(models.Book{ID: 1, Title: "Go Concurrency", Author: "John Doe"})
	library.AddBook(models.Book{ID: 2, Title: "Parallel Programming", Author: "Alice"})

	controller := &controllers.LibraryController{
		Service:  library,
		Requests: make(chan concurrency.ReservationRequest, 10),
	}

	controller.Start()
}
