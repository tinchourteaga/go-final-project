package domain

type ReportCarries struct {
	LocalityID   string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarriesCount int    `json:"carries_count"`
}
