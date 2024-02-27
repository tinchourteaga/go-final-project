package requests

type EmployeeDTOPost struct {
	CardNumberID *string `json:"card_number_id" binding:"required"`
	FirstName    *string `json:"first_name" binding:"required"`
	LastName     *string `json:"last_name" binding:"required"`
	WarehouseID  *int    `json:"warehouse_id" binding:"required"`
}

type EmployeeDTOPatch struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	WarehouseID int    `json:"warehouse_id"`
}
