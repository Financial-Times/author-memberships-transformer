package main

type berthaRole struct {
	UUID       string `json:"uuid"`
	Preflabel  string `json:"preflabel"`
	ParentUuid string `json:"parentUuid,omitempty"`
}
