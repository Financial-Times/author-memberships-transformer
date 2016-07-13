package main

type membership struct {
	UUID                   string                 `json:"uuid"`
	PrefLabel              string                 `json:"prefLabel,omitempty"`
	PersonUUID             string                 `json:"personUuid"`
	OrganisationUUID       string                 `json:"organisationUuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers"`
	MembershipRoles        []membershipRole       `json:"membershipRoles"`
}

type alternativeIdentifiers struct {
	UUIDS []string `json:"uuids"`
}

type membershipRole struct {
	RoleUUID string `json:"roleuuid,omitempty"`
}
