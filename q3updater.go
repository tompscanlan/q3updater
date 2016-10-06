package q3updater

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Verbose lets other packages know we want to be noisey
var Verbose = true

// Team is the team id
var Team int

type JournalEntry struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}


type SoapReservation struct {
	Name       string    `json:"name,omitempty"`
	StartDate  time.Time `json:"start_date,omitempty"`
	EndDate    time.Time `json:"end_date,omitempty"`
	ServerName string    `json:"server_name,omitempty"`
	Approved   bool      `json:"approved"`
}

type Reservation struct {
	Approved *bool `json:"approved"`
	Begin time.Time `json:"begin"`
	End time.Time `json:"end"`
	Username *string `json:"username"`
}

var (
	JournalId = 0

	JournalServer  string
	ApprovalServer string
	LabDataServer  string
)

func NewJournalEntryFromJson(jsonStr []byte) *JournalEntry {
	log.Println("in NewJournalEntryFromJson")

	entry := new(JournalEntry)
	err := json.Unmarshal(jsonStr, entry)
	if err != nil {
		log.Println(err)
	}

	return entry
}

func NewReservationFromJournalEntry(entry JournalEntry) *Reservation {
	log.Println("in NewReservationFromJournalEntry")

	res := new(Reservation)

	decoded, err := base64.StdEncoding.DecodeString(entry.Message)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(decoded, res)
	if err != nil {
		log.Println(err)
	}

	return res
}


func (r Reservation) Bytes() []byte {
	b, err := json.Marshal(r)

	if err != nil {
		log.Println(err)
	}
	return b
}
func (r Reservation) String() string {

	b := r.Bytes()
	s := string(b[:])
	return s
}

func (r SoapReservation) Bytes() []byte {
	b, err := json.Marshal(r)

	if err != nil {
		log.Println(err)
	}
	return b
}
func (r SoapReservation) String() string {

	b := r.Bytes()
	s := string(b[:])
	return s
}

func GetJournalEntry(server string) (*JournalEntry, error) {
	log.Println("GetJournalEntry: Getting journal entry")
	entry := new(JournalEntry)

	// create the request
	url := fmt.Sprintf("%s/%s", server, "api/topic/10")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("GetJournalEntry:", err)
		return entry, err
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("GetJournalEntry:", err)
		return entry, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if Verbose {
		log.Println("GetJournalEntry: response Status:", resp.Status)
		log.Println("GetJournalEntry: response Headers:", resp.Header)

		log.Println("GetJournalEntry: response Body:", string(body))
	}

	err = json.Unmarshal(body, entry)
	if err != nil {
		log.Println("GetJournalEntry:", err)
		return entry, err
	}

	return entry, nil
}

func PostApprovedToReservation(server string, approved *Approval) error {

	log.Println("PostApprovedToReservation: sending approved to reservation")

	if !approved.Approved {
		return errors.New("Attempted to register a non-approved reservation")
	}

	decoded, err := base64.StdEncoding.DecodeString(approved.Description)
	if err != nil {
		log.Println(err)
		return err
	}

	// convert from a soapUI object
	reservation := new(SoapReservation)
	err = json.Unmarshal(decoded, reservation)
	if err != nil {
		log.Println(err)
		return err
	}
		newReservation := new(Reservation)

	newReservation.Username = &reservation.Name
	newReservation.Begin = reservation.StartDate
	newReservation.End = reservation.EndDate
	newReservation.Approved = &approved.Approved

	if reservation.Name == "" {
		return errors.New("decoded invalid reservation")
	}

	// create the request
	url := fmt.Sprintf("%s/%s/%s/%s", server, "item", reservation.ServerName, "reservation")
	log.Printf("PostApprovedToReservation: %s, data: %s", url, newReservation)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(newReservation.Bytes()))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if Verbose {
		log.Println("PostApprovedToReservation: response Status:", resp.Status)
		log.Println("PostApprovedToReservation: response Headers:", resp.Header)

		log.Println("PostApprovedToReservation: response Body:", string(body))
	}
	return nil
}
