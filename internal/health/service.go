package health

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Live(ctx context.Context) HealthResponse {
	return HealthResponse{
		Status: "ok",
	}
}

func (s *Service) Ready(ctx context.Context) HealthResponse {
	return HealthResponse{
		Status: "ok",
		Checks: map[string]string{
			"postgres": "not_configured",
			"redis":    "not_configured",
			"kafka":    "not_configured",
		},
	}
}
