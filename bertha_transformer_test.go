package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var anAuthorUuid = "1eaeabeb-01ef-42dc-bf95-459e9087a763"
var aRoleLabel = "Superhero"
var anotherRoleLabel = "Rockstar"
var aJobTitle = "Avengers member"
var aMembershipUuid = "0311bc2c-1dc1-489b-a202-59252c28d2de"
var aRoleUuid = "b4f06685-9f58-40af-850f-07f1585fab73"

var anAuthor = author{
	UUID:           anAuthorUuid,
	Role:           aRoleLabel,
	Jobtitle:       aJobTitle,
	Membershipuuid: aMembershipUuid,
}

var anotherAuthor = author{
	UUID:           "1eaeabeb-01ef-42dc-bf95-459e9087a763",
	Role:           anotherRoleLabel,
	Jobtitle:       aJobTitle,
	Membershipuuid: aMembershipUuid,
}

var aBerthaRole = berthaRole{
	UUID:      aRoleUuid,
	Preflabel: aRoleLabel,
}

var aRolesMap = map[string]berthaRole{aBerthaRole.Preflabel: aBerthaRole}

var aMembership = membership{
	UUID:                   aMembershipUuid,
	PrefLabel:              aJobTitle,
	PersonUUID:             anAuthorUuid,
	OrganisationUUID:       ftUuid,
	AlternativeIdentifiers: alternativeIdentifiers{UUIDS: []string{aMembershipUuid}},
	MembershipRoles:        []membershipRole{membershipRole{RoleUUID: aRoleUuid}},
}

func TestShouldTransformAuthorToPersonSucessfully(t *testing.T) {
	transformer := berthaTransformer{}
	m, err := transformer.toMembership(anAuthor, aRolesMap)
	assert.Nil(t, err)
	assert.Equal(t, aMembership, m, "The membership is transformed properly")
}

func TestShouldReturnErrorWhenRoleIsNotFound(t *testing.T) {
	transformer := berthaTransformer{}
	_, err := transformer.toMembership(anotherAuthor, aRolesMap)
	assert.NotNil(t, err)
}
