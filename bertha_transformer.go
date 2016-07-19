package main

import (
	"errors"
	"fmt"
)

const ftUuid = "dac01f07-4b6d-3615-8532-a56752cc7e5f"

type berthaTransformer struct {
}

func (bt *berthaTransformer) toMembership(a author, uuidRolesMap map[string]berthaRole, namesRolesMap map[string]berthaRole) (membership, error) {

	memRoles, err := bt.buildMembershipRoles(a.Role, uuidRolesMap, namesRolesMap)

	if err != nil {
		return membership{}, err
	}

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

func (bt *berthaTransformer) buildMembershipRoles(roleName string, uuidRolesMap map[string]berthaRole, nameRolesMap map[string]berthaRole) ([]membershipRole, error) {
	berthaRole := nameRolesMap[roleName]
	memRoles := []membershipRole{}
	if berthaRole.UUID == "" {
		return []membershipRole{}, errors.New(fmt.Sprintf(`Role UUID is not found for "%s"`, berthaRole.Preflabel))
	}

	for parentRoleUuid := berthaRole.UUID; parentRoleUuid != ""; {
		berthaRole = uuidRolesMap[parentRoleUuid]
		memRole, err := bt.transformRole(berthaRole)
		if err != nil {
			return []membershipRole{}, err
		}
		memRoles = append(memRoles, memRole)
		parentRoleUuid = berthaRole.ParentUuid
	}
	return memRoles, nil
}

func (bt *berthaTransformer) transformRole(br berthaRole) (membershipRole, error) {
	if br.UUID == "" {
		return membershipRole{}, errors.New(fmt.Sprintf(`Role UUID is not found for "%s"`, br.Preflabel))
	}
	return membershipRole{RoleUUID: br.UUID}, nil
}
