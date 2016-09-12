package q3updater

import (
	"encoding/json"
	"log"
	"time"
)

// Verbose lets other packages know we want to be noisey

var Verbose = true

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
