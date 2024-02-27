package requests

type InboundOrderDTOPOST struct {
	OrderDate      *string `json:"order_date" binding:"required"`
	OrderNumber    *string `json:"order_number" binding:"required"`
	EmployeeID     *int    `json:"employee_id" binding:"required"`
	ProductBatchID *int    `json:"product_batch_id" binding:"required"`
	WarehouseID    *int    `json:"warehouse_id" binding:"required"`
}
