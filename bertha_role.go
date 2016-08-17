package main

type berthaRole struct {
	UUID       string `json:"uuid"`
	Preflabel  string `json:"preflabel"`
	ParentUUID string `json:"parentUuid,omitempty"`
}
