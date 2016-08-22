package main

import (
	"fmt"

	"github.com/pborman/uuid"
)

const ftUUID = "dac01f07-4b6d-3615-8532-a56752cc7e5f"

type berthaTransformer struct {
}

func (bt *berthaTransformer) toMembership(a author, uuidRolesMap map[string]berthaRole, namesRolesMap map[string]berthaRole) (membership, error) {

	memRoles, err := bt.buildMembershipRoles(a.Role, uuidRolesMap, namesRolesMap)

	if err != nil {
		return membership{}, err
	}

	personUUID := uuid.NewMD5(uuid.UUID{}, []byte(a.TmeIdentifier)).String()

	membershipUUID := uuid.NewMD5(uuid.UUID{}, []byte(personUUID+"_MEMBER_"+ftUUID)).String()

	altIds := alternativeIdentifiers{
		UUIDS: []string{membershipUUID},
	}

	m := membership{
		UUID:                   membershipUUID,
		PrefLabel:              a.Jobtitle,
		PersonUUID:             personUUID,
		OrganisationUUID:       ftUUID,
		AlternativeIdentifiers: altIds,
		MembershipRoles:        memRoles,
	}
	return m, nil
}

func (bt *berthaTransformer) buildMembershipRoles(roleName string, uuidRolesMap map[string]berthaRole, nameRolesMap map[string]berthaRole) ([]membershipRole, error) {
	berthaRole := nameRolesMap[roleName]
	memRoles := []membershipRole{}
	if berthaRole.UUID == "" {
		return []membershipRole{}, fmt.Errorf(`Role UUID is not found for "%s"`, berthaRole.Preflabel)
	}

	for parentRoleUUID := berthaRole.UUID; parentRoleUUID != ""; {
		berthaRole = uuidRolesMap[parentRoleUUID]
		memRole, err := bt.transformRole(berthaRole)
		if err != nil {
			return []membershipRole{}, err
		}
		memRoles = append(memRoles, memRole)
		parentRoleUUID = berthaRole.ParentUUID
	}
	return memRoles, nil
}

func (bt *berthaTransformer) transformRole(br berthaRole) (membershipRole, error) {
	if br.UUID == "" {
		return membershipRole{}, fmt.Errorf(`Role UUID is not found for "%s"`, br.Preflabel)
	}
	return membershipRole{RoleUUID: br.UUID}, nil
}
