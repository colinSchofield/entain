package service

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/colinSchofield/entain/sporting/db"
	"github.com/colinSchofield/entain/sporting/logging"
	"github.com/colinSchofield/entain/sporting/proto/sporting"
)

type Sporting interface {
	// ListSports will return a collection of sports.
	ListSports(ctx context.Context, in *sporting.ListSportsRequest) (*sporting.ListSportsResponse, error)
	// Return a sport based upon its id
	GetSport(ctx context.Context, in *sporting.GetSportRequest) (*sporting.GetSportResponse, error)
}

// sportingService implements the Sporting interface.
type sportingService struct {
	sportsRepo db.SportsRepo
}

// NewSportingService instantiates and returns a new sportingService.
func NewSportingService(sportsRepo db.SportsRepo) Sporting {
	return &sportingService{sportsRepo}
}

func (s *sportingService) ListSports(ctx context.Context, in *sporting.ListSportsRequest) (*sporting.ListSportsResponse, error) {
	sports, err := s.sportsRepo.List(in.Filter, in.Order)
	if err != nil {
		wrappedError := fmt.Errorf("unexpected error occurred in call to Repo List: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	}

	logging.Logger().Debugf("%d sports were returned to the caller", len(sports))
	return &sporting.ListSportsResponse{Sports: sports}, nil
}

func (s *sportingService) GetSport(ctx context.Context, in *sporting.GetSportRequest) (*sporting.GetSportResponse, error) {
	if sport, err := s.sportsRepo.Get(in.Id); err != nil {
		logging.Logger().Errorf("Unexpected error message %s", err)
		return nil, err
	} else {
		return &sporting.GetSportResponse{Sport: sport}, nil
	}
}
