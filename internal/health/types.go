package health

type HealthResponse struct{
	Status string `json:"status"`
	Checks map[string] string `json:checks,omitempty`
}