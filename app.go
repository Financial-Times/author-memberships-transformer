package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	app := cli.App("author-memberhips-transformer", "A RESTful API for transforming Bertha Curated Authors to UP People Memberships JSON")

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
		bs := newBerthaService(*berthaAuthorsSrcUrl, *berthaRolesSrcUrl)
		mh := newMembershipHandler(bs)

		h := setupServiceHandlers(mh)

		http.Handle("/", httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
			httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), h)))

		log.Infof("Listening on [%d].", *port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		if err != nil {
			log.Printf("Web server failed: [%v].", err)
		}
	}

	app.Run(os.Args)
}

func setupServiceHandlers(mh membershipHandler) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc(status.PingPath, status.PingHandler)
	r.HandleFunc(status.PingPathDW, status.PingHandler)
	r.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
	r.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
	r.HandleFunc("/__health", v1a.Handler("Curated Authors Membership Transformer", "Checks for accessing Bertha", mh.AuthorsHealthCheck(), mh.RolesHealthCheck()))
	r.HandleFunc(status.GTGPath, mh.GoodToGo)

	r.HandleFunc("/transformers/author-memberships/__count", mh.getMembershipsCount).Methods("GET")
	r.HandleFunc("/transformers/author-memberships/__ids", mh.getMembershipUuids).Methods("GET")
	r.HandleFunc("/transformers/author-memberships/{uuid}", mh.getMembershipByUuid).Methods("GET")

	return r
}
