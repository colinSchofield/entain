package service

import (
	"context"
	"errors"
	"testing"

	"github.com/colinSchofield/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
)

// Mock returning an empty value
type racesRepoMock struct{}

func (r *racesRepoMock) Init() error {
	return nil
}
func (r *racesRepoMock) List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestOrderBy) ([]*racing.Race, error) {
	return nil, nil
}

func (r *racesRepoMock) Get(id int64) (*racing.Race, error) {
	return nil, nil
}

// Mock returning a database error
type racesRepoMockError struct{}

func (r *racesRepoMockError) Init() error {
	return errors.New("Db error")
}

func (r *racesRepoMockError) List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestOrderBy) ([]*racing.Race, error) {
	return nil, errors.New("Db error")
}

func (r *racesRepoMockError) Get(id int64) (*racing.Race, error) {
	return nil, errors.New("Db error")
}

func Test_ListRacesReturnsAnEmptyValue(t *testing.T) {
	// Given
	racingService := NewRacingService(&racesRepoMock{})
	request := new(racing.ListRacesRequest)
	ctx := context.TODO()
	// When
	results, err := racingService.ListRaces(ctx, request)
	// Then
	assert.NotNil(t, results)
	assert.Zero(t, len(results.Races), "Mock returns an empty value")
	assert.Nil(t, err)
}

func Test_ListRacesReturnsAnError(t *testing.T) {
	// Given
	racingService := NewRacingService(&racesRepoMockError{})
	request := new(racing.ListRacesRequest)
	ctx := context.TODO()
	// When
	results, err := racingService.ListRaces(ctx, request)
	// Then
	assert.Nil(t, results)
	assert.NotNil(t, err)
}

func Test_GetRacesReturnsAnEmptyValue(t *testing.T) {
	// Given
	racingService := NewRacingService(&racesRepoMock{})
	request := new(racing.GetRaceRequest)
	ctx := context.TODO()
	// When
	results, err := racingService.GetRace(ctx, request)
	// Then
	assert.NotNil(t, results)
	assert.Nil(t, results.Race, "Mock returns an empty value")
	assert.Nil(t, err)
}

func Test_GetRaceReturnsAnError(t *testing.T) {
	// Given
	racingService := NewRacingService(&racesRepoMockError{})
	request := new(racing.GetRaceRequest)
	ctx := context.TODO()
	// When
	results, err := racingService.GetRace(ctx, request)
	// Then
	assert.Nil(t, results)
	assert.NotNil(t, err)
}
