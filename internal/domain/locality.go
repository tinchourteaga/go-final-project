package domain

type Locality struct {
	ID           string `json:"id" binding:"required"`
	LocalityName string `json:"locality_name" binding:"required"`
	ProvinceName string `json:"province_name" binding:"required"`
	CountryName  string `json:"country_name" binding:"required"`
}
