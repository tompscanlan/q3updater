package main

import (
	//	"encoding/json"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/tompscanlan/labreserved/client"
	//	"github.com/tompscanlan/labreserved/models"

	//	"github.com/tompscanlan/labreserved/client/operations"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/tompscanlan/q3updater"
	"log"
	"net/http"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	port    = kingpin.Flag("port", "port to listen on").Default(listenPortDefault).OverrideDefaultFromEnvar("PORT").Short('l').Int()

	journalServer  = kingpin.Flag("journal-server", "REST endpoint for the journal server").Default(journalServerDefault).OverrideDefaultFromEnvar("JOURNAL_SERVER").Short('j').String()
	approvalServer = kingpin.Flag("approval-server", "REST endpoint for the approval server").Default(approvalServerDefault).OverrideDefaultFromEnvar("APPROVAL_SERVER").Short('a').String()
	labDataServer  = kingpin.Flag("labdata-server", "REST endpoint for the lab data server").Default(labDataServerDefault).OverrideDefaultFromEnvar("LABDATA_SERVER").Short('d').String()
	lock           = sync.RWMutex{}

	Client        *apiclient.Labreserved
	AllActive     = q3updater.NewActive()
	JournalTicker = time.NewTicker(time.Second * 1)
)

const (
	listenPortDefault     = "8083"
	journalServerDefault  = "http://journal.butterhead.net:8080"
	approvalServerDefault = "http://q3.butterhead.net:2080"
	labDataServerDefault  = "labreserved.butterhead.net:2080"
)

func init() {
	setupFlags()
	q3updater.Verbose = *verbose
}

func setupFlags() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()
}

func main() {

	transport := httptransport.New(*labDataServer, "", []string{"http"})
	Client = apiclient.New(transport, strfmt.Default)

	api := rest.NewApi()

	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		// record hit stats
		rest.Get("/.status", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(statusMw.GetStatus())
		}),

		rest.Get("/active", GetActive),
		rest.Put("/active", PutActive),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	CheckJournal()
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), api.MakeHandler())
	JournalTicker.Stop()
	fmt.Println("Ticker stopped")
	log.Fatal(err)
}

func CheckJournal() {
	go func() {
		for t := range JournalTicker.C {
			fmt.Println("Tick at", t)
			je, err := q3updater.GetJournalEntry(*journalServer)

			if err != nil {
				log.Println(err)
			}
			res := q3updater.NewReservationFromJournalEntry(*je)
			_ = res
		}
	}()
}

func GetActive(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	w.WriteJson(AllActive)
	lock.RUnlock()
}

func PutActive(w rest.ResponseWriter, r *rest.Request) {
	active := q3updater.NewActive()

	var b []byte
	_, err := r.Body.Read(b)
	if err != nil {
		log.Println(err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		log.Println("body: ", string(b[:]))
	}

	err = r.DecodeJsonPayload(&active)
	if err != nil {
		log.Println(err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("active: %s", active.String())

	lock.RLock()
	AllActive.Active = active.Active

	err = w.WriteJson(active)
	if err != nil {
		log.Println(err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lock.RUnlock()
}
