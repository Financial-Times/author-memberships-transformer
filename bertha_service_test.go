package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const etag = "W/\"75e-78600296\""
const aurhorsBerthaPath = "/view/publish/gss/123456XYZ/Authors"
const rolesBerthaPath = "/view/publish/gss/123456XYZ/Roles"
const authorsBerthaOutput = "test-resources/bertha-authors-output.json"
const rolesBerthaOutput = "test-resources/bertha-roles-output.json"

var membership1 = membership{
	UUID:                   "e6e8b382-4833-11e6-beb8-9e71128cae77",
	PrefLabel:              "Chief Economics Commentator",
	PersonUUID:             "daf5fed2-013c-468d-85c4-aee779b8aa53",
	OrganisationUUID:       "dac01f07-4b6d-3615-8532-a56752cc7e5f",
	AlternativeIdentifiers: alternativeIdentifiers{UUIDS: []string{"e6e8b382-4833-11e6-beb8-9e71128cae77"}},
	MembershipRoles:        []membershipRole{membershipRole{RoleUUID: "7ef75a6a-b6bf-4eb7-a1da-03e0acabef1b"}},
}
var membership2 = membership{
	UUID: "c721a241-9250-4c77-9620-8abb08027686",
}

type berthaMock struct {
	server     *httptest.Server
	outputFile string
	path       string
}

var berthaAuthorsMock = berthaMock{outputFile: authorsBerthaOutput, path: aurhorsBerthaPath}
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
	uuids := bs.getMembershipUuids()
	assert.Equal(t, 2, len(uuids), "Bertha should return 2 authors")
	assert.Equal(t, true, contains(uuids, membership1.UUID), "It should contain membership1 UUID")
	assert.Equal(t, true, contains(uuids, membership2.UUID), "It should contain membership2 UUID")
}

func TestShouldReturnSingleMembership(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	bs.getMembershipCount()
	m := bs.getMembershipByUuid(membership1.UUID)
	assert.Equal(t, membership1, m, "The membership should be membership1")
}

func TestShouldReturnEmptyMembershipUuidsWhenMembershipCountIsNotCalled(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	uuids := bs.getMembershipUuids()
	assert.Equal(t, 0, len(uuids), "Bertha should return 0 memberships")
}

func TestShouldReturnEmptyMembershiprWhenMembershipCountIsNotCalled(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	m := bs.getMembershipByUuid(membership1.UUID)
	assert.Equal(t, membership{}, m, "The membership should be empty")
}

func TestShouldReturnEmptyMembershipWhenMembershipIsNotAvailable(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	m := bs.getMembershipByUuid("7f8bd61a-3575-4d32-a758-0fa41cbcc826")
	assert.Equal(t, membership{}, m, "The membership should be empty")
}

func TestShouldReturnErrorWhenBerthaAuthorsIsUnhappy(t *testing.T) {
	berthaAuthorsMock.start("unhappy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("happy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c, err := bs.getMembershipCount()
	assert.NotNil(t, err)
	assert.Equal(t, -1, c, "It should return -1")

	uuids := bs.getMembershipUuids()
	assert.Equal(t, 0, len(uuids), "It should return 0 UUIDs")

	m := bs.getMembershipByUuid(membership1.UUID)
	assert.Equal(t, membership{}, m, "The author should be empty")
}

func TestShouldReturnErrorWhenBerthaRolesIsUnhappy(t *testing.T) {
	berthaAuthorsMock.start("happy")
	defer berthaAuthorsMock.stop()
	berthaRolesMock.start("unhappy")
	defer berthaRolesMock.stop()

	bs := newBerthaService(berthaAuthorsMock.getUrl(), berthaRolesMock.getUrl())

	c, err := bs.getMembershipCount()
	assert.NotNil(t, err)
	assert.Equal(t, -1, c, "It should return -1")

	uuids := bs.getMembershipUuids()
	assert.Equal(t, 0, len(uuids), "It should return 0 UUIDs")

	m := bs.getMembershipByUuid(membership1.UUID)
	assert.Equal(t, membership{}, m, "The author should be empty")
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
