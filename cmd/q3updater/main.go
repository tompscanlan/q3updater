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
	ticker  = kingpin.Flag("ticker", "seconds between polling events").Default(fmt.Sprintf("%d", tickerDefault)).Short('t').Int64()
	team    = kingpin.Flag("team", "team id").Default(fmt.Sprintf("%d", teamDefault)).OverrideDefaultFromEnvar("TEAM_ID").Int()

	journalServer  = kingpin.Flag("journal-server", "REST endpoint for the journal server").Default(journalServerDefault).OverrideDefaultFromEnvar("JOURNAL_SERVER").Short('j').String()
	approvalServer = kingpin.Flag("approval-server", "REST endpoint for the approval server").Default(approvalServerDefault).OverrideDefaultFromEnvar("APPROVAL_SERVER").Short('a').String()
	labDataServer  = kingpin.Flag("labdata-server", "REST endpoint for the lab data server").Default(labDataServerDefault).OverrideDefaultFromEnvar("LABDATA_SERVER").Short('d').String()
	lock           = sync.RWMutex{}

	Client         *apiclient.Labreserved
	journalTicker  *time.Ticker
	approvalTicker *time.Ticker
)

const (
	teamDefault           = 7357
	tickerDefault         = 5
	listenPortDefault     = "8083"
	journalServerDefault  = "http://journal.butterhead.net:8080"
	approvalServerDefault = "http://approval.vmwaredevops.appspot.com"
	labDataServerDefault  = "labreserved.butterhead.net:2080"
)

func init() {
	setupFlags()
	q3updater.Verbose = *verbose
	q3updater.Team = *team
	q3updater.JournalServer = *journalServer
	q3updater.ApprovalServer = *approvalServer
	q3updater.LabDataServer = *labDataServer

	q3updater.JournalTicker = time.NewTicker(time.Second * time.Duration(*ticker))
	q3updater.ApprovalTicker = time.NewTicker(time.Second * time.Duration(*ticker))

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

		rest.Get("/active", q3updater.GetActive),
		rest.Put("/active", q3updater.PutActive),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	// Continually check the Journal for new reservations
	// and send them on for approval
	q3updater.SendReservationForApproval()

	// also poll approval service and send approved
	// on to be recorded
	q3updater.SendApprovedToReserved()

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), api.MakeHandler())
	q3updater.JournalTicker.Stop()
	q3updater.ApprovalTicker.Stop()

	fmt.Println("Ticker stopped")
	log.Fatal(err)
}
