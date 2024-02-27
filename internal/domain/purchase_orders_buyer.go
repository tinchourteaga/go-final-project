package domain

type Purchase_orders_buyer struct {
	ID           int    `json:"buyer_id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	OrdersCount  int    `json:"orders_count"`
}
