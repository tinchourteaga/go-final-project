package requests

type SellerPostRequest struct {
	CID         *int    `json:"cid" binding:"required"`
	CompanyName *string `json:"company_name" binding:"required"`
	Address     *string `json:"address" binding:"required"`
	Telephone   *string `json:"telephone" binding:"required"`
	Locality_id *string `json:"locality_id" binding:"required"`
}

type SellerPatchRequest struct {
	CID         *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
	Locality_id *string `json:"locality_id"`
}
