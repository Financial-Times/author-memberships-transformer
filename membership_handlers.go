package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type membershipHandler struct {
	membershipService membershipService
}

func newMembershipHandler(ms membershipService) membershipHandler {
	return membershipHandler{
		membershipService: ms,
	}
}

func (mh *membershipHandler) getMembershipsCount(writer http.ResponseWriter, req *http.Request) {
	c, err := mh.membershipService.getMembershipCount()
	if err != nil {
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
	} else {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf(`%v`, c))
		buffer.WriteTo(writer)
	}
}

func (mh *membershipHandler) getMembershipUuids(writer http.ResponseWriter, req *http.Request) {
	uuids := mh.membershipService.getMembershipUuids()
	writeStreamResponse(uuids, writer)
}

func (mh *membershipHandler) getMembershipByUuid(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	m := mh.membershipService.getMembershipByUuid(uuid)
	writeJSONResponse(m, !reflect.DeepEqual(m, membership{}), writer)
}

func (mh *membershipHandler) AuthorsHealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for curated author data from Bertha",
		Name:             "Check connectivity to Bertha",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/curated-authors-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to Bertha to be able to supply curated authors",
		Checker:          mh.authorsChecker,
	}
}

func (mh *membershipHandler) authorsChecker() (string, error) {
	err := mh.membershipService.checkAuthorsConnectivity()
	if err == nil {
		return "Connectivity to Bertha Authors is ok", err
	}
	return "Error connecting to Bertha Authors", err
}

func (mh *membershipHandler) RolesHealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for curated author roles data from Bertha",
		Name:             "Check connectivity to Bertha",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/curated-authors-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to Bertha to be able to supply curated authors",
		Checker:          mh.rolesChecker,
	}
}

func (mh *membershipHandler) rolesChecker() (string, error) {
	err := mh.membershipService.checkRolesConnectivity()
	if err == nil {
		return "Connectivity to Bertha Authors is ok", err
	}
	return "Error connecting to Bertha Authors", err
}

func (mh *membershipHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := mh.authorsChecker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}

func writeStreamResponse(ids []string, writer http.ResponseWriter) {
	for _, id := range ids {
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf(`{"id":"%s"} `, id))
		buffer.WriteTo(writer)
	}
}
