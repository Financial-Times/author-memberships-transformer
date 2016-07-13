package main

type membershipService interface {
	getMembershipCount() (int, error)
	getMembershipUuids() []string
	getMembershipByUuid(uuid string) membership
	checkAuthorsConnectivity() error
	checkRolesConnectivity() error
}
