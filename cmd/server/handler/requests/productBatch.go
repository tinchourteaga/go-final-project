package requests

type PostProductBatch struct {
	BatchNumber        int    `json:"batch_number" binding:"required"`
	CurrentQuantity    int    `json:"current_quantity" binding:"required"`
	CurrentTemperature int    `json:"current_temperature" binding:"required"`
	DueDate            string `json:"due_date" binding:"required"`
	InitialQuantity    int    `json:"initial_quantity" binding:"required"`
	ManufacturingDate  string `json:"manufacturing_date" binding:"required"`
	ManufacturingHour  int    `json:"manufacturing_hour" binding:"required"`
	MinimumTemperature int    `json:"minimum_temperature" binding:"required"`
	ProductID          int    `json:"product_id" binding:"required"`
	SectionID          int    `json:"section_id" binding:"required"`
}
