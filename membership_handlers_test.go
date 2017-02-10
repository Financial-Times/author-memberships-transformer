package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var curatedAuthorsMembershipTransformer *httptest.Server

//For fixtures see fixtures_test.go

var uuids = []string{expectedMembershipUUID, "e06be0f8-0426-4ee8-80e3-3da37255818a"}
var expectedStreamOutput = "{\"id\":\"" + expectedMembershipUUID + "\"}\n{\"id\":\"e06be0f8-0426-4ee8-80e3-3da37255818a\"}\n"

type MockedBerthaService struct {
	mock.Mock
}

func (m *MockedBerthaService) refreshMembershipCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedBerthaService) getMembershipUuids() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockedBerthaService) getMembershipByUuid(uuid string) membership {
	args := m.Called(uuid)
	return args.Get(0).(membership)
}

func (m *MockedBerthaService) getMembershipCount() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockedBerthaService) checkAuthorsConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedBerthaService) checkRolesConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

func startCuratedAuthorsMembershipTransformer(bs *MockedBerthaService) {
	mh := newMembershipHandler(bs)
	h := setupServiceHandlers(mh)
	curatedAuthorsMembershipTransformer = httptest.NewServer(h)
}

func TestShouldReturn200AndMembershipCount(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipCount").Return(2)
	mbs.On("refreshMembershipCache").Return(nil)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/memberships/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be text/plain")
	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, "2", actualOutput, "Response body should contain the count of available authors")
}

func TestShouldReturn500WhenMembershipCountIsCalledAndCacheRefreshFails(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipCount").Return(2)
	mbs.On("refreshMembershipCache").Return(errors.New("Exterminate!"))
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/memberships/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status should be 500")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type should be application/json")
	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, "{\"message\": \"Exterminate!\"}\n", actualOutput, "Response body should contain the error message")
}

func TestShouldReturn200WhenMembershipCacheIsRefreshed(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("refreshMembershipCache").Return(nil)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Post(curatedAuthorsMembershipTransformer.URL+"/transformers/memberships/__reload", "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldReturn200AndMembershipUuids(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipUuids").Return(uuids, nil)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/memberships/__ids")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be text/plain")
	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, expectedStreamOutput, actualOutput, "Response body should be a sequence of ids")
}

func getStringFromReader(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

func TestShouldReturn200AndTransformedMembership(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipByUuid", expectedMembershipUUID).Return(expectedMembership, nil)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/memberships/" + expectedMembershipUUID)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "Content-Type should be application/json")

	file, _ := os.Open("test-resources/transformed-membership-output.json")
	defer file.Close()

	expectedOutput := getStringFromReader(file)
	actualOutput := getStringFromReader(resp.Body)

	assert.JSONEq(t, expectedOutput, actualOutput, "Response body should be a valid membership")
}

func TestShouldReturn404WhenMembershipIsNotFound(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipByUuid", expectedMembershipUUID).Return(membership{}, nil)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/memberships/" + expectedMembershipUUID)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response status should be 404")
}

func TestShouldReturn500WhenCacheRefreshReturnsError(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("refreshMembershipCache").Return(errors.New("I am a zombie"))
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Post(curatedAuthorsMembershipTransformer.URL+"/transformers/memberships/__reload", "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status should be 500")
}
