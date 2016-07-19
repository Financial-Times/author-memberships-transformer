package main

type transformer interface {
	toMembership(author, map[string]berthaRole, map[string]berthaRole) (membership, error)
}
