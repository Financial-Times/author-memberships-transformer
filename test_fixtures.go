package main

var anAuthorTmeIdentifier = "Q0ItMDAwMDkwMA==-QXV0aG9ycw=="

// MD5 hash of the author TME identifier
var expectedAuthorUUID = "0f07d468-fc37-3c44-bf19-a81f2aae9f36"

// MD5 hash using both the expected author uuid and organisation uuid
var expectedMembershipUUID = "78a23be4-b7b0-392a-a900-582a0dbe383b"

var aRoleLabel = "Superhero"
var anotherRoleLabel = "Rockstar"
var aJobTitle = "Avengers member"

var aRoleUUID = "b4f06685-9f58-40af-850f-07f1585fab73"
var anotherRoleUUID = "b4f06685-9f58-40af-850f-07f1585fab73"
var yetAnotherRoleUUID = "93e60bde-dc80-4ed3-8cc1-a19346c52014"

var anotherAuthorTmeIdentifier = "Q0ItMDAwMDkyNg==-QXV0aG9ycw=="

// MD5 hash of the anotherAuthor TME Identifier
var expectedAnotherAuthorUUID = "8f9ac45f-2cc2-35f7-83f4-579c66a09eb0"

var anAuthor = author{
	Role:          aRoleLabel,
	Jobtitle:      aJobTitle,
	TmeIdentifier: anAuthorTmeIdentifier,
}

var anotherAuthor = author{
	Role:          anotherRoleLabel,
	Jobtitle:      aJobTitle,
	TmeIdentifier: "",
}

var aBerthaRole = berthaRole{
	UUID:       aRoleUUID,
	Preflabel:  aRoleLabel,
	ParentUUID: yetAnotherRoleUUID,
}

var anotherBerthaRole = berthaRole{
	UUID:      yetAnotherRoleUUID,
	Preflabel: "Hero",
}

var aNameRolesMap = map[string]berthaRole{aBerthaRole.Preflabel: aBerthaRole, anotherBerthaRole.Preflabel: anotherBerthaRole}
var aUUIDRolesMap = map[string]berthaRole{aBerthaRole.UUID: aBerthaRole, anotherBerthaRole.UUID: anotherBerthaRole}

var expectedMembership = membership{
	UUID:                   expectedMembershipUUID,
	PrefLabel:              aJobTitle,
	PersonUUID:             expectedAuthorUUID,
	OrganisationUUID:       ftUUID,
	AlternativeIdentifiers: alternativeIdentifiers{UUIDS: []string{expectedMembershipUUID}},
	MembershipRoles:        []membershipRole{membershipRole{RoleUUID: aRoleUUID}, membershipRole{RoleUUID: yetAnotherRoleUUID}},
}
