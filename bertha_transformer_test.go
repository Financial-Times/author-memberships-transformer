package main

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var anAuthorUuid = uuid.NewV4()
var aRoleLabel = "Superhero"
var anotherRoleLabel = "Rockstar"
var aJobTitle = "Avengers member"
var aMembershipUuid = uuid.NewV4()
var aRoleUuid = uuid.NewV4()

var anAuthor = author{
	UUID:           anAuthorUuid,
	Role:           aRoleLabel,
	Jobtitle:       aJobTitle,
	Membershipuuid: aMembershipUuid,
}

var anotherAuthor = author{
	UUID:           uuid.NewV4(),
	Role:           anotherRoleLabel,
	Jobtitle:       aJobTitle,
	Membershipuuid: aMembershipUuid,
}

var aBerthaRole = bertaRole{
	UUID:       aRoleUuid,
	aRoleLabel: aRoleLabel,
}

var aRoleLabel = map[string]berthaRole{aRoleLabel.aRoleLabel: aBerthaRole}

var aMembership = membership{
	UUID:                   aMembershipUuid,
	PrefLabel:              aJobTitle,
	PersonUUID:             anAuthorUuid,
	OrganisationUUID:       ftUuid,
	AlternativeIdentifiers: alternativeIdentifiers{UUIDS: []string{aMembershipUuid}},
	MembershipRoles:        []membershipRole{membershipRole{RoleUUID: aRoleUuids}},
}

func TestShouldTransformAuthorToPersonSucessfully(t *testing.T) {
	transformer := berthaTransformer{}
	m, err := transformer.toMembership(anAuthor, aRolesMap)
	assert.Nil(t, err)
	assert.Equal(t, aMembership, m, "The membership is transformed properly")
}

func TestShouldReturnErrorWhenRoleIsNotFound(t *testing.T) {
	transformer := berthaTransformer{}
	m, err := transformer.toMembership(anotherAuthor, aRolesMap)
	assert.NotNil(t, err)
}
