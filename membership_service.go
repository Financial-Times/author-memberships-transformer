package main

type membershipService interface {
	refreshMembershipCache() error
	getMembershipCount() int
	getMembershipUuids() []string
	getMembershipByUuid(uuid string) membership
	checkAuthorsConnectivity() error
	checkRolesConnectivity() error
}
