package requests

type WarehousePostRequest struct {
	Address            *string `json:"address" binding:"required"`
	Telephone          *string `json:"telephone" binding:"required"`
	WarehouseCode      *string `json:"warehouse_code" binding:"required"`
	MinimumCapacity    *int    `json:"minimum_capacity" binding:"required"`
	MinimumTemperature *int    `json:"minimum_temperature" binding:"required"`
}

type WarehousePatchRequest struct {
	Address            *string `json:"address"`
	Telephone          *string `json:"telephone"`
	WarehouseCode      *string `json:"warehouse_code"`
	MinimumCapacity    *int    `json:"minimum_capacity"`
	MinimumTemperature *int    `json:"minimum_temperature"`
}
