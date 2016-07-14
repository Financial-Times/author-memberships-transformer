package main

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var curatedAuthorsMembershipTransformer *httptest.Server

var uuids = []string{aMembership.UUID, "e06be0f8-0426-4ee8-80e3-3da37255818a"}
var expectedStreamOutput = `{"id":"` + aMembership.UUID + `"} {"id":"e06be0f8-0426-4ee8-80e3-3da37255818a"} `

type MockedBerthaService struct {
	mock.Mock
}

func (m *MockedBerthaService) getMembershipUuids() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockedBerthaService) getMembershipByUuid(uuid string) membership {
	args := m.Called(uuid)
	return args.Get(0).(membership)
}

func (m *MockedBerthaService) getMembershipCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
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
	mbs.On("getMembershipCount").Return(2, nil)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/author-memberships/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be text/plain")
	actualOutput := getStringFromReader(resp.Body)
	assert.Equal(t, "2", actualOutput, "Response body should contain the count of available authors")
}

func TestShouldReturn200AndMembershipUuids(t *testing.T) {

	mbs := new(MockedBerthaService)
	mbs.On("getMembershipUuids").Return(uuids)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/author-memberships/__ids")
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

func TestShouldReturn200AndTrasformedMembership(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipByUuid", aMembership.UUID).Return(aMembership)
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/author-memberships/" + aMembership.UUID)
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

	assert.Equal(t, expectedOutput, actualOutput, "Response body should be a valid membership")
}

func TestShouldReturn404WhenMembershipIsNotFound(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipByUuid", aMembership.UUID).Return(membership{})
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/author-memberships/" + aMembership.UUID)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response status should be 404")
}

func TestShouldReturn500WhenBerthaReturnsError(t *testing.T) {
	mbs := new(MockedBerthaService)
	mbs.On("getMembershipCount").Return(-1, errors.New("I am a zobie"))
	startCuratedAuthorsMembershipTransformer(mbs)
	defer curatedAuthorsMembershipTransformer.Close()

	resp, err := http.Get(curatedAuthorsMembershipTransformer.URL + "/transformers/author-memberships/__count")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Response status should be 500")
}
