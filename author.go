package main

// This struct reflects the JSON data model of curated authors from Bertha
type author struct {
	UUID           string `json:"uuid"`
	Role           string `json:"role"`
	Jobtitle       string `json:"jobtitle"`
	Membershipuuid string `json:"membershipuuid"`
}
