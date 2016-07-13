package main

type transformer interface {
	toMembership(author, map[string]berthaRole) (membership, error)
}
