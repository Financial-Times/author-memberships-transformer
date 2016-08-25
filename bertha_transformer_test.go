package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//For fixtures see fixtures_test.go

func TestShouldTransformAuthorToPersonSuccessfully(t *testing.T) {
	transformer := berthaTransformer{}
	m, err := transformer.toMembership(anAuthor, aUUIDRolesMap, aNameRolesMap)
	assert.Nil(t, err)
	assert.Equal(t, expectedMembership, m, "The membership is transformed properly")
}

func TestShouldReturnErrorWhenRoleIsNotFound(t *testing.T) {
	transformer := berthaTransformer{}
	_, err := transformer.toMembership(anotherAuthor, aUUIDRolesMap, aNameRolesMap)
	assert.NotNil(t, err)
}
