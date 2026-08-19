package service

import (
	"github.com/chandankrr/loreline/internal/lib/job"
	"github.com/chandankrr/loreline/internal/repository"
	"github.com/chandankrr/loreline/internal/server"
)

type Services struct {
	Job *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	return &Services{
		Job: s.Job,
	}, nil
}
