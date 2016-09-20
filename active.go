package q3updater

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sync"
	"time"
)

type Active struct {
	Id     int  `json:"id,omitempty"`
	Active bool `json:"active,bool"`
}

var JournalTicker *time.Ticker
var ApprovalTicker *time.Ticker

var Activity = NewActive()
var lock = sync.RWMutex{}

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

func GetActive(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	w.WriteJson(Activity)
	lock.RUnlock()
}

func PutActive(w rest.ResponseWriter, r *rest.Request) {
	active := NewActive()

	err := r.DecodeJsonPayload(&active)
	if err != nil {
		log.Println(err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("active: %s", active.String())

	lock.RLock()
	Activity.Active = active.Active
	Activity.Id = active.Id + 1
	lock.RUnlock()

	w.WriteJson(active)
}
