package domain

// Product represents a table of fresh products on the database.
type Product struct {
	ID                             int     `json:"id"`
	Description                    string  `json:"description"`
	ExpirationRate                 int     `json:"expiration_rate"`
	FreezingRate                   int     `json:"freezing_rate"`
	Height                         float32 `json:"height"`
	Length                         float32 `json:"length"`
	NetWeight                      float32 `json:"net_weight"`
	ProductCode                    string  `json:"product_code"`
	RecommendedFreezingTemperature float32 `json:"recommended_freezing_temperature"`
	Width                          float32 `json:"width"`
	ProductTypeID                  int     `json:"product_type_id"`
	SellerID                       *int    `json:"seller_id,omitempty"`
}
