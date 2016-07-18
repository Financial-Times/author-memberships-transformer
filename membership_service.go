package main

type membershipService interface {
	refreshMembershipCache() error
	getMembershipCount() (int, error)
	getMembershipUuids() ([]string, error)
	getMembershipByUuid(uuid string) (membership, error)
	checkAuthorsConnectivity() error
	checkRolesConnectivity() error
}
