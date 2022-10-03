package db

import (
	"testing"
	"time"

	"github.com/colinSchofield/entain/racing/proto/racing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_HappyPath(t *testing.T) {
	// Given
	filterArgs := []int64{5, 6}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
		AddRow("1", "5", "North Dakota foes", "3", false, time.Now())
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE meeting_id IN (.+)").
		WithArgs(filterArgs[0], filterArgs[1]).
		WillReturnRows(rows)
	defer db.Close()
	racesRepo := NewRacesRepo(db)
	filter := &racing.ListRacesRequestFilter{MeetingIds: filterArgs}
	// When
	if raceResults, err := racesRepo.List(filter); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, raceResults)
		assert.True(t, raceResults[0].MeetingId == int64(5))
	}
}

func Test_EmptyFilter(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
		AddRow("1", "2", "North Dakota foes", "3", false, time.Now())
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races").
		WillReturnRows(rows)
	defer db.Close()
	racesRepo := NewRacesRepo(db)
	filter := &racing.ListRacesRequestFilter{}
	// When
	if raceResults, err := racesRepo.List(filter); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, raceResults)
	}
}

func Test_NoRowsReturned(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"})
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races").
		WillReturnRows(rows)
	defer db.Close()
	racesRepo := NewRacesRepo(db)
	filter := &racing.ListRacesRequestFilter{}
	// When
	if raceResults, err := racesRepo.List(filter); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.Nil(t, raceResults)
	}
}

func Test_VisibleFilterTrue(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
		AddRow("1", "2", "North Dakota foes", "3", true, time.Now())
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE visible = 1").
		WithArgs().
		WillReturnRows(rows)
	defer db.Close()
	racesRepo := NewRacesRepo(db)
	filter := &racing.ListRacesRequestFilter{Visible: true}
	// When
	if raceResults, err := racesRepo.List(filter); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, raceResults)
		assert.True(t, raceResults[0].Visible)
	}
}

func Test_BothMeetingIdAndVisibleFilters(t *testing.T) {
	// Given
	filterArgs := []int64{5, 6}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
		AddRow("1", "5", "North Dakota foes", "3", true, time.Now())
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races WHERE meeting_id IN (.+) AND visible = 1").
		WithArgs(filterArgs[0], filterArgs[1]).
		WillReturnRows(rows)
	defer db.Close()
	racesRepo := NewRacesRepo(db)
	filter := &racing.ListRacesRequestFilter{Visible: true, MeetingIds: filterArgs}
	// When
	if raceResults, err := racesRepo.List(filter); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, raceResults)
		assert.True(t, raceResults[0].MeetingId == int64(5))
	}
}

func Test_VisibleSetFalseFilter(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
		AddRow("1", "2", "North Dakota foes", "3", false, time.Now())
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM races").
		WithArgs().
		WillReturnRows(rows)
	defer db.Close()
	racesRepo := NewRacesRepo(db)
	filter := &racing.ListRacesRequestFilter{Visible: true}
	// When
	if raceResults, err := racesRepo.List(filter); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, raceResults)
		assert.False(t, raceResults[0].Visible)
	}
}
