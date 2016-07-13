package main

import (
	"errors"
	"fmt"
)

const ftUuid = "dac01f07-4b6d-3615-8532-a56752cc7e5f"

type berthaTransformer struct {
}

func (bt *berthaTransformer) toMembership(a author, rolesMap map[string]berthaRole) (membership, error) {
	fmt.Println(a.Membershipuuid)

	roleUuid := rolesMap[a.Role].UUID

	if roleUuid == "" {
		return membership{}, errors.New(fmt.Sprintf(`Role UUID is not found for "%s"`, a.Role))
	}

	memRole := membershipRole{RoleUUID: roleUuid}
	memRoles := []membershipRole{memRole}

	altIds := alternativeIdentifiers{
		UUIDS: []string{a.Membershipuuid},
	}

	m := membership{
		UUID:                   a.Membershipuuid,
		PrefLabel:              a.Jobtitle,
		PersonUUID:             a.UUID,
		OrganisationUUID:       ftUuid,
		AlternativeIdentifiers: altIds,
		MembershipRoles:        memRoles,
	}
	return m, nil
}
