package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("curated-authors-memberships-transformer", "A RESTful API for transforming Bertha Curated Authors to UP People Memberships JSON")

	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})
	berthaAuthorsSrcUrl := app.String(cli.StringOpt{
		Name:   "bertha-authors-source-url",
		Value:  "{url}",
		Desc:   "The URL of the Bertha Authors JSON source",
		EnvVar: "BERTHA_AUTHORS_SOURCE_URL",
	})
	berthaRolesSrcUrl := app.String(cli.StringOpt{
		Name:   "bertha-roles-source-url",
		Value:  "{url}",
		Desc:   "The URL of the Bertha Roles JSON source",
		EnvVar: "BERTHA_ROLES_SOURCE_URL",
	})

	app.Action = func() {
		log.Info("App started!!!")
		bs, err := newBerthaService(*berthaAuthorsSrcUrl, *berthaRolesSrcUrl)

		if err != nil {
			log.Error(err)
			panic(err)
		}

		mh := newMembershipHandler(bs)

		h := setupServiceHandlers(mh)

		http.Handle("/", httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
			httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), h)))

		log.Infof("Listening on [%d].", *port)
		errServe := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		if errServe != nil {
			log.Printf("Web server failed: [%v].", errServe)
		}
	}

	app.Run(os.Args)
}

func setupServiceHandlers(mh membershipHandler) http.Handler {
	r := mux.NewRouter()

	timedHC := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  "curated-authors-memberships-tf",
			Name:        "Curated Authors Memberships Transformer",
			Description: "A REST service that transforms Authors data from Bertha to Memberships according to UPP format.",
			Checks:      []fthealth.Check{mh.AuthorsHealthCheck(), mh.RolesHealthCheck()},
		},
		Timeout: 10 * time.Second,
	}

	r.HandleFunc("/__health", fthealth.Handler(timedHC))
	r.HandleFunc(status.PingPath, status.PingHandler)
	r.HandleFunc(status.PingPathDW, status.PingHandler)
	r.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
	r.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
	r.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(mh.GTG))

	r.HandleFunc("/transformers/memberships/__reload", mh.refreshMembershipCache).Methods("POST")
	r.HandleFunc("/transformers/memberships/__count", mh.getMembershipsCount).Methods("GET")
	r.HandleFunc("/transformers/memberships/__ids", mh.getMembershipUuids).Methods("GET")
	r.HandleFunc("/transformers/memberships/{uuid}", mh.getMembershipByUuid).Methods("GET")

	return r
}
