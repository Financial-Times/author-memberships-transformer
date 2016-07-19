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
var anotherRoleUuid = "b4f06685-9f58-40af-850f-07f1585fab73"
var yetAnotherRoleUuid = "93e60bde-dc80-4ed3-8cc1-a19346c52014"

var anAuthor = author{
	UUID:           anAuthorUuid,
	Role:           aRoleLabel,
	Jobtitle:       aJobTitle,
	Membershipuuid: aMembershipUuid,
}

var anotherAuthor = author{
	UUID:           "57d23eb1-65dc-4f67-ba4b-f76ca60efccc",
	Role:           anotherRoleLabel,
	Jobtitle:       aJobTitle,
	Membershipuuid: aMembershipUuid,
}

var aBerthaRole = berthaRole{
	UUID:       aRoleUuid,
	Preflabel:  aRoleLabel,
	ParentUuid: yetAnotherRoleUuid,
}

var anotherBerthaRole = berthaRole{
	UUID:      yetAnotherRoleUuid,
	Preflabel: "Hero",
}

var aNameRolesMap = map[string]berthaRole{aBerthaRole.Preflabel: aBerthaRole, anotherBerthaRole.Preflabel: anotherBerthaRole}
var aUuidRolesMap = map[string]berthaRole{aBerthaRole.UUID: aBerthaRole, anotherBerthaRole.UUID: anotherBerthaRole}

var aMembership = membership{
	UUID:                   aMembershipUuid,
	PrefLabel:              aJobTitle,
	PersonUUID:             anAuthorUuid,
	OrganisationUUID:       ftUuid,
	AlternativeIdentifiers: alternativeIdentifiers{UUIDS: []string{aMembershipUuid}},
	MembershipRoles:        []membershipRole{membershipRole{RoleUUID: aRoleUuid}, membershipRole{RoleUUID: yetAnotherRoleUuid}},
}

func TestShouldTransformAuthorToPersonSucessfully(t *testing.T) {
	transformer := berthaTransformer{}
	m, err := transformer.toMembership(anAuthor, aUuidRolesMap, aNameRolesMap)
	assert.Nil(t, err)
	assert.Equal(t, aMembership, m, "The membership is transformed properly")
}

func TestShouldReturnErrorWhenRoleIsNotFound(t *testing.T) {
	transformer := berthaTransformer{}
	_, err := transformer.toMembership(anotherAuthor, aUuidRolesMap, aNameRolesMap)
	assert.NotNil(t, err)
}
