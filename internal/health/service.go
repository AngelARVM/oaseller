package health

import "context"

type Checker interface {
	Check(ctx context.Context) error
}

type Service struct {
	postgres Checker
}

func NewService(postgres Checker) *Service {
	return &Service{
		postgres: postgres,
	}
}

func (s *Service) Live(ctx context.Context) HealthResponse {
	return HealthResponse{
		Status: "ok",
	}
}

func (s *Service) Ready(ctx context.Context) HealthResponse {
	checks := map[string]string{
		"postgres": "ok",
		"redis":    "not_configured",
		"kafka":    "not_configured",
	}

	status := "ok"

	if err := s.postgres.Check(ctx); err != nil {
		checks["postgres"] = "error"
		status = "error"
	}

	return HealthResponse{
		Status: status,
		Checks: checks,
	}
}
