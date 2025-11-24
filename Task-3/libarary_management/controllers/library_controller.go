package controllers

import (
	"fmt"
	"library_management/concurrency"
	"library_management/services"
	"time"
)

type LibraryController struct {
	Service  *services.Library
	Requests chan concurrency.ReservationRequest
}

func (lc *LibraryController) Start() {
	fmt.Println("Starting concurrent reservation simulation...")

	// Start worker
	concurrency.StartReservationWorker(lc.Service, lc.Requests)

	// Simulate multiple concurrent members reserving same book
	go func() {
		lc.Requests <- concurrency.ReservationRequest{BookID: 1, MemberID: 101}
	}()
	go func() {
		lc.Requests <- concurrency.ReservationRequest{BookID: 1, MemberID: 102}
	}()
	go func() {
		lc.Requests <- concurrency.ReservationRequest{BookID: 2, MemberID: 103}
	}()

	time.Sleep(6 * time.Second)
	fmt.Println("Reservation simulation complete.")
}
