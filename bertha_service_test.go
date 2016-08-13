package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const etag = "W/\"75e-78600296\""
const authorsBerthaPath = "/view/publish/gss/123456XYZ/Authors"
const rolesBerthaPath = "/view/publish/gss/123456XYZ/Roles"
const authorsBerthaOutput = "test-resources/bertha-authors-output.json"
const rolesBerthaOutput = "test-resources/bertha-roles-output.json"

var membership1 = membership{
	UUID:                   expectedMembershipUUID,
	PrefLabel:              "Chief Economics Commentator",
	PersonUUID:             expectedAuthorUUID,
	OrganisationUUID:       ftUUID,
	AlternativeIdentifiers: alternativeIdentifiers{UUIDS: []string{expectedMembershipUUID}},
	MembershipRoles:        []membershipRole{membershipRole{RoleUUID: "7ef75a6a-b6bf-4eb7-a1da-03e0acabef1b"}},
}
var membership2 = membership{
	UUID: "a1c08d1f-9c19-370b-af34-80aa6cf3c0ad",
}

type berthaMock struct {
	server     *httptest.Server
	outputFile string
	path       string
}

var berthaAuthorsMock = berthaMock{outputFile: authorsBerthaOutput, path: authorsBerthaPath}
var berthaRolesMock = berthaMock{outputFile: rolesBerthaOutput, path: rolesBerthaPath}

func (mock *berthaMock) getUrl() string {
	return mock.server.URL + mock.path
}

func (mock *berthaMock) start(status string) {
	r := mux.NewRouter()
	if status == "happy" {
		r.Path(mock.path).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(mock.berthaHandlerMock)})
	} else {
		r.Path(mock.path).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(unhappyHandler)})
	}
	mock.server = httptest.NewServer(r)
}

func (mock *berthaMock) berthaHandlerMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ifNoneMatch := r.Header.Get("If-None-Match")

	if ifNoneMatch == etag {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.Header().Set("ETag", etag)

		file, err := os.Open(mock.outputFile)
		if err != nil {
			return
		}
		defer file.Close()
		io.Copy(w, file)
	}
}

func (mock *berthaMock) stop() {
	mock.server.Close()
}

func unhappyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TestShouldReturnMembershipCount(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c, err := bs.getMembershipCount()

	assert.Nil(t, err)
	assert.Equal(t, 2, c, "Bertha should return 2 authors")
}

func TestShouldReturnMembershipsUuids(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	bs.getMembershipCount()
	uuids, err := bs.getMembershipUuids()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(uuids), "Bertha should return 2 authors")
	assert.Equal(t, true, contains(uuids, membership1.UUID), "actual UUIDS=%s should contain expected membership1 UUID=%s", uuids, membership1.UUID)
	assert.Equal(t, true, contains(uuids, membership2.UUID), "actual UUIDS=%s should contain expected membership2 UUID=%s", uuids, membership2.UUID)
}

func TestShouldReturnSingleMembership(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	bs.getMembershipCount()
	m, err := bs.getMembershipByUuid(membership1.UUID)

	assert.Nil(t, err)
	assert.Equal(t, membership1, m, "The membership should be membership1")
}

func TestShouldReturnEmptyMembershipWhenMembershipIsNotAvailable(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())
	m, err := bs.getMembershipByUuid("7f8bd61a-3575-4d32-a758-0fa41cbcc826")

	assert.Nil(t, err)
	assert.Equal(t, membership{}, m, "The membership should be empty")
}

func TestShouldReturnErrorWhenBerthaAuthorsIsUnhappy(t *testing.T) {
	berthaAuthorsMock.start("unhappy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	err := bs.refreshMembershipCache()
	assert.NotNil(t, err)

	c, err := bs.getMembershipCount()
	assert.NotNil(t, err)
	assert.Equal(t, -1, c, "It should return -1")

	uuids, err := bs.getMembershipUuids()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(uuids), "It should return 0 UUIDs")

	m, err := bs.getMembershipByUuid(membership1.UUID)
	assert.NotNil(t, err)
	assert.Equal(t, membership{}, m, "The membership should be empty")
}

func TestShouldReturnErrorWhenBerthaRolesIsUnhappy(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("unhappy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	err := bs.refreshMembershipCache()
	assert.NotNil(t, err)

	c, err := bs.getMembershipCount()
	assert.NotNil(t, err)
	assert.Equal(t, -1, c, "It should return -1")

	uuids, err := bs.getMembershipUuids()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(uuids), "It should return 0 UUIDs")

	m, err := bs.getMembershipByUuid(membership1.UUID)
	assert.NotNil(t, err)
	assert.Equal(t, membership{}, m, "The membership should be empty")
}

func TestCheckConnectivityOfHappyBertaAuthors(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c := bs.checkAuthorsConnectivity()
	assert.Nil(t, c)
}

func TestCheckConnectivityOfUnhappyBerthaAuthors(t *testing.T) {
	berthaAuthorsMock.start("unhappy")
	defer berthaAuthorsMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c := bs.checkAuthorsConnectivity()
	assert.NotNil(t, c)
}

func TestCheckConnectivityOfHappyBertaRoles(t *testing.T) {
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c := bs.checkRolesConnectivity()
	assert.Nil(t, c)
}

func TestCheckConnectivityOfUnhappyBerthaRoles(t *testing.T) {
	berthaRolesMock.start("unhappy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c := bs.checkRolesConnectivity()
	assert.NotNil(t, c)
}

func TestCheckConnectivityBerthaAuthorsOffline(t *testing.T) {
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c := bs.checkAuthorsConnectivity()
	assert.NotNil(t, c)
}

func TestCheckConnectivityBerthaRolesOffline(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c := bs.checkRolesConnectivity()
	assert.NotNil(t, c)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
