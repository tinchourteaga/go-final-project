package domain

type Employee struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

type EmployeeWithInboundOrders struct {
	Employee
	InboundOrders int `json:"inbound_orders_count"`
}
