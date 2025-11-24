package concurrency

import (
	"fmt"
	"library_management/services"
)

type ReservationRequest struct {
	BookID   int
	MemberID int
}

func StartReservationWorker(library *services.Library, requests chan ReservationRequest) {
	go func() {
		for req := range requests {
			err := library.ReserveBook(req.BookID, req.MemberID)
			if err != nil {
				fmt.Printf("[ERROR] Reservation failed for Book %d by Member %d: %s\n",
					req.BookID, req.MemberID, err.Error())
			} else {
				fmt.Printf("[SUCCESS] Book %d reserved by Member %d\n", req.BookID, req.MemberID)
			}
		}
	}()
}
