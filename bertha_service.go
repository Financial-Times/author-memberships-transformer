package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gregjones/httpcache"
	"net/http"
)

var client = httpcache.NewMemoryCacheTransport().Client()

type berthaService struct {
	authorsUrl     string
	rolesUrl       string
	membershipsMap map[string]membership
	transformer    transformer
}

func newBerthaService(authorsUrl string, rolesUrl string) *berthaService {
	return &berthaService{
		authorsUrl:  authorsUrl,
		rolesUrl:    rolesUrl,
		transformer: &berthaTransformer{},
	}
}

func (bs *berthaService) getMembershipCount() (int, error) {
	bs.membershipsMap = make(map[string]membership)
	authResp, authErr := bs.callBerthaService(bs.authorsUrl)
	if authErr != nil {
		log.Error(authErr)
		return -1, authErr
	}

	var authors []author
	if err := json.NewDecoder(authResp.Body).Decode(&authors); err != nil {
		log.Error(err)
		return -1, err
	}

	rolesResp, rolesErr := bs.callBerthaService(bs.rolesUrl)
	if rolesErr != nil {
		log.Error(rolesErr)
		return -1, rolesErr
	}

	var roles []berthaRole
	if err := json.NewDecoder(rolesResp.Body).Decode(&roles); err != nil {
		log.Error(err)
		return -1, err
	}

	if err := bs.populateMembershipMap(authors, roles); err != nil {
		log.Error(err)
		return -1, err
	}

	return len(bs.membershipsMap), nil
}

func (bs *berthaService) populateMembershipMap(authors []author, roles []berthaRole) error {
	rolesMap := make(map[string]berthaRole)

	for _, r := range roles {
		rolesMap[r.Preflabel] = r
	}

	for _, a := range authors {
		m, err := bs.transformer.toMembership(a, rolesMap)
		if err != nil {
			bs.membershipsMap = make(map[string]membership)
			return err
		}
		bs.membershipsMap[m.UUID] = m
	}
	return nil
}

func (bs *berthaService) getMembershipUuids() []string {
	uuids := make([]string, 0)
	for uuid, _ := range bs.membershipsMap {
		uuids = append(uuids, uuid)
	}
	return uuids
}

func (bs *berthaService) getMembershipByUuid(uuid string) membership {
	return bs.membershipsMap[uuid]
}

func (bs *berthaService) callBerthaService(url string) (res *http.Response, err error) {
	log.WithFields(log.Fields{"bertha_url": url}).Info("Calling Bertha...")
	res, err = client.Get(url)
	return
}

func (bs *berthaService) checkAuthorsConnectivity() error {
	return bs.checkConnectivity(bs.authorsUrl)
}

func (bs *berthaService) checkRolesConnectivity() error {
	return bs.checkConnectivity(bs.rolesUrl)
}

func (bs *berthaService) checkConnectivity(url string) error {
	resp, err := bs.callBerthaService(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Bertha returns unexpected HTTP status: %d", resp.StatusCode))
	}
	return nil
}
