package service

import (
	"context"
	"errors"
	"testing"

	"github.com/colinSchofield/entain/sporting/proto/sporting"
	"github.com/stretchr/testify/assert"
)

// Mock returning an empty value
type sportsRepoMock struct{}

func (r *sportsRepoMock) Init() error {
	return nil
}
func (r *sportsRepoMock) List(filter *sporting.ListSportsRequestFilter, sort *sporting.ListSportsRequestOrderBy) ([]*sporting.Sport, error) {
	return nil, nil
}

func (r *sportsRepoMock) Get(id int64) (*sporting.Sport, error) {
	return nil, nil
}

// Mock returning a database error
type sportsRepoMockError struct{}

func (r *sportsRepoMockError) Init() error {
	return errors.New("Db error")
}

func (r *sportsRepoMockError) List(filter *sporting.ListSportsRequestFilter, sort *sporting.ListSportsRequestOrderBy) ([]*sporting.Sport, error) {
	return nil, errors.New("Db error")
}

func (r *sportsRepoMockError) Get(id int64) (*sporting.Sport, error) {
	return nil, errors.New("Db error")
}

func Test_ListSportsReturnsAnEmptyValue(t *testing.T) {
	// Given
	sportingService := NewSportingService(&sportsRepoMock{})
	request := new(sporting.ListSportsRequest)
	ctx := context.TODO()
	// When
	results, err := sportingService.ListSports(ctx, request)
	// Then
	assert.NotNil(t, results)
	assert.Zero(t, len(results.Sports), "Mock returns an empty value")
	assert.Nil(t, err)
}

func Test_ListSportsReturnsAnError(t *testing.T) {
	// Given
	sportingService := NewSportingService(&sportsRepoMockError{})
	request := new(sporting.ListSportsRequest)
	ctx := context.TODO()
	// When
	results, err := sportingService.ListSports(ctx, request)
	// Then
	assert.Nil(t, results)
	assert.NotNil(t, err)
}

func Test_GetSportsReturnsAnEmptyValue(t *testing.T) {
	// Given
	sportingService := NewSportingService(&sportsRepoMock{})
	request := new(sporting.GetSportRequest)
	ctx := context.TODO()
	// When
	results, err := sportingService.GetSport(ctx, request)
	// Then
	assert.NotNil(t, results)
	assert.Nil(t, results.Sport, "Mock returns an empty value")
	assert.Nil(t, err)
}

func Test_GetSportReturnsAnError(t *testing.T) {
	// Given
	sportingService := NewSportingService(&sportsRepoMockError{})
	request := new(sporting.GetSportRequest)
	ctx := context.TODO()
	// When
	results, err := sportingService.GetSport(ctx, request)
	// Then
	assert.Nil(t, results)
	assert.NotNil(t, err)
}
