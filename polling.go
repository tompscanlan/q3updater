package q3updater

import (
	"log"
)

var LastApprovedSent = 0

func SendReservationForApproval() {
	go func() {
		for t := range JournalTicker.C {
			log.Printf("Send Reservation Tick at %v", t)

			// skip this check if we're in active
			if !Activity.Active {
				log.Println("Skipping tick")
				continue
			}

			// get a journal entry
			je, err := GetJournalEntry(JournalServer)
			if err != nil {
				log.Println(err)
				continue
			}
			// send it for approval
			if je.Message != "" {
				log.Printf("got journal Entry and sending it to ApprovalServer: %v", *je)

				err := PostForApproval(ApprovalServer, je)
							if err != nil {

				log.Printf("failed to send a journal entry to approval: %s", err)
				}
			}

		}
	}()
}

///TODO.... {"name":"testres-servername1","start_date":"2016-01-01T06:00:00Z","end_date":"2099-12-31T00:00:00Z","server_name":"servername1","approved":false}
//// start_date/end_date should follow models.reservations
func SendApprovedToReserved() {
	go func() {
		for t := range ApprovalTicker.C {
			log.Printf("Record Approved Tick at %v", t)

			// skip this check if we're in active
			if !Activity.Active {
				log.Println("Skipping tick")
				continue
			}

			approved, err := GetApproved(ApprovalServer, Team)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, a := range approved {

				if a.Description == "" || a.Approved == false || LastApprovedSent >= a.Id {
					continue
				}
				log.Printf("Found %d approvals. operating on %d", len(approved), a.Id)

				err := PostApprovedToReservation(LabDataServer, &a)
				if err != nil {
					log.Printf("failed to send an approved reservation: %s", err)
					continue
				}
				LastApprovedSent = a.Id

			}

		}
	}()
}
