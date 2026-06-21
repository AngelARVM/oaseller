package merchants

type CreateMerchantRequest struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
