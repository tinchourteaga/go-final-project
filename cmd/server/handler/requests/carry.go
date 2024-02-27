package requests

type CarryPostRequest struct {
	CID         *string `json:"cid" binding:"required"`
	CompanyName *string `json:"company_name" binding:"required"`
	Address     *string `json:"address" binding:"required"`
	Telephone   *string `json:"telephone" binding:"required"`
	Locality_id *string `json:"locality_id" binding:"required"`
}
