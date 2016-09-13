package q3updater

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Verbose lets other packages know we want to be noisey
var Verbose = true

type JournalEntry struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type Reservation struct {
	Name       string    `json:"name,omitempty"`
	StartDate  time.Time `json:"start_date,omitempty"`
	EndDate    time.Time `json:"end_date,omitempty"`
	ServerName string    `json:"server_name,omitempty"`
}

type Approval struct {
	Id          int    `json:"id,omitempty"`
	TeamId      int    `json:"teamID,omitempty"`
	Blob        int    `json:"blob,omitempty"`
	Description string `json:"description,omitempty"`
	Approved    bool   `json:"approved,omitempty"`
}

type Active struct {
	Id     int  `json:"id,omitempty"`
	Active bool `json:"active,bool"`
}

func NewActive() Active {
	a := new(Active)
	a.Enable()
	return *a
}

func (a *Active) Enable() {

	a.Id = 1
	a.Active = true
}

func (a *Active) Disable() {

	a.Id = 1
	a.Active = true
}

// returns JSON string
func (a Active) String() string {
	b, err := json.Marshal(a)

	if err != nil {
		log.Println(err)
	}

	s := string(b[:])
	return s
}

var JournalId = 0

func NewJournalEntryFromJson(jsonStr []byte) *JournalEntry {
	entry := new(JournalEntry)
	err := json.Unmarshal(jsonStr, entry)
	if err != nil {
		log.Println(err)
	}

	return entry
}

func NewReservationFromJournalEntry(entry JournalEntry) *Reservation {
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

func GetJournalEntry(server string) (*JournalEntry, error) {
	log.Println("Getting journal entry")
	entry := new(JournalEntry)

	// create the request
	url := fmt.Sprintf("%s/%s", server, "api/topic/10")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return entry, err
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return entry, err
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	return entry, nil
}
