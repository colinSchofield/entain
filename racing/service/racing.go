package service

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/colinSchofield/entain/racing/db"
	"github.com/colinSchofield/entain/racing/logging"
	"github.com/colinSchofield/entain/racing/proto/racing"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter, in.Order)
	if err != nil {
		wrappedError := fmt.Errorf("unexpected error occurred in call to Repo List: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	}

	logging.Logger().Debugf("%d races were returned to the caller", len(races))
	return &racing.ListRacesResponse{Races: races}, nil
}
