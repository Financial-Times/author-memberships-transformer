package main

// This struct reflects the JSON data model of curated authors from Bertha
type author struct {
	Role          string `json:"role"`
	Jobtitle      string `json:"jobtitle"`
	TmeIdentifier string `json:"tmeidentifier"`
}
