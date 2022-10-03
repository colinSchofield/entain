package db

import (
	"testing"
	"time"

	"github.com/colinSchofield/entain/sporting/proto/sporting"

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
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports WHERE meeting_id IN (.+)").
		WithArgs(filterArgs[0], filterArgs[1]).
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	filter := &sporting.ListSportsRequestFilter{MeetingIds: filterArgs}
	order := &sporting.ListSportsRequestOrderBy{}
	// When
	if sportResults, err := sportsRepo.List(filter, order); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, sportResults)
		assert.True(t, sportResults[0].MeetingId == int64(5))
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
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports").
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	filter := &sporting.ListSportsRequestFilter{}
	order := &sporting.ListSportsRequestOrderBy{}
	// When
	if sportResults, err := sportsRepo.List(filter, order); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, sportResults)
	}
}

func Test_NoRowsReturned(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"})
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports").
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	filter := &sporting.ListSportsRequestFilter{}
	order := &sporting.ListSportsRequestOrderBy{}
	// When
	if sportResults, err := sportsRepo.List(filter, order); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.Nil(t, sportResults)
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
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports WHERE visible = 1").
		WithArgs().
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	filter := &sporting.ListSportsRequestFilter{Visible: true}
	order := &sporting.ListSportsRequestOrderBy{}
	// When
	if sportResults, err := sportsRepo.List(filter, order); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, sportResults)
		assert.True(t, sportResults[0].Visible)
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
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports WHERE meeting_id IN (.+) AND visible = 1").
		WithArgs(filterArgs[0], filterArgs[1]).
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	filter := &sporting.ListSportsRequestFilter{Visible: true, MeetingIds: filterArgs}
	order := &sporting.ListSportsRequestOrderBy{}
	// When
	if sportResults, err := sportsRepo.List(filter, order); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, sportResults)
		assert.True(t, sportResults[0].MeetingId == int64(5))
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
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports").
		WithArgs().
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	filter := &sporting.ListSportsRequestFilter{Visible: true}
	order := &sporting.ListSportsRequestOrderBy{}
	// When
	if sportResults, err := sportsRepo.List(filter, order); err != nil {
		t.Errorf("List returned an error of: %s", err)
	} else {
		// Then
		assert.NotNil(t, sportResults)
		assert.False(t, sportResults[0].Visible)
	}
}

func Test_OrderByMultiTest(t *testing.T) {

	multiTest := []struct {
		scenario        string
		attributeName   string
		directionClause string
		expectedQuery   string
	}{
		{scenario: "Happy Path Ascending", attributeName: "advertisedStartTime", directionClause: "ASC", expectedQuery: " ORDER BY advertised_start_time ASC"},
		{scenario: "Happy Path Descending", attributeName: "advertisedStartTime", directionClause: "DESC", expectedQuery: " ORDER BY advertised_start_time DESC"},
		{scenario: "Happy Path Descending with different attribute name", attributeName: "name", directionClause: "ASC", expectedQuery: " ORDER BY name ASC"},
		{scenario: "SQL Injection test", attributeName: "SQL Injection!!", directionClause: "ASC", expectedQuery: " ORDER BY advertised_start_time ASC"},
	}

	for _, test := range multiTest {
		// Given
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"})
		mock.ExpectQuery(test.expectedQuery).
			WillReturnRows(rows)
		defer db.Close()
		sportsRepo := NewSportsRepo(db)
		var direction sporting.SortDirection
		if test.directionClause == "ASC" {
			direction = sporting.SortDirection_ASC
		} else {
			direction = sporting.SortDirection_DESC
		}
		filter := &sporting.ListSportsRequestFilter{}
		order := &sporting.ListSportsRequestOrderBy{OrderBy: test.attributeName, Direction: direction}
		// When
		sportResults, err := sportsRepo.List(filter, order)
		// Then
		assert.Nil(t, sportResults, test.scenario)
		assert.Nil(t, err, test.scenario)

	}
}

func Test_SportStatusMultiTest(t *testing.T) {

	multiTest := []struct {
		scenario    string
		rowTime     time.Time
		sportStatus string
	}{
		{scenario: "Sport far in the Future", rowTime: time.Now().AddDate(10, 0, 0), sportStatus: "OPEN"},
		{scenario: "Sport far in the Past", rowTime: time.Now().AddDate(-10, 0, 0), sportStatus: "CLOSED"},
		{scenario: "Sport has JUST closed", rowTime: time.Now().AddDate(0, 0, 0), sportStatus: "CLOSED"},
	}

	for _, test := range multiTest {
		// Given
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
			AddRow("1", "2", "North Dakota foes", "3", false, test.rowTime)
		mock.ExpectQuery(" ").
			WillReturnRows(rows)
		defer db.Close()
		sportsRepo := NewSportsRepo(db)
		filter := &sporting.ListSportsRequestFilter{}
		order := &sporting.ListSportsRequestOrderBy{}
		// When
		sportResults, err := sportsRepo.List(filter, order)
		// Then
		assert.Equal(t, sportResults[0].GetStatus().String(), test.sportStatus, test.scenario)
		assert.Nil(t, err, test.scenario)
	}
}

func Test_GetRequestHappyPath(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"}).
		AddRow("1", "2", "North Dakota foes", "3", false, time.Now())
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports").
		WithArgs().
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	// When
	sport, err := sportsRepo.Get(1)
	// Then
	assert.Nil(t, err)
	assert.NotNil(t, sport)
}

func Test_GetRequestNotFound(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "meeting_id", "name", "number", "visible", "advertised_start_time"})
	mock.ExpectQuery("SELECT id, meeting_id, name, number, visible, advertised_start_time FROM sports").
		WithArgs().
		WillReturnRows(rows)
	defer db.Close()
	sportsRepo := NewSportsRepo(db)
	// When
	sport, err := sportsRepo.Get(1)
	// Then
	assert.NotNil(t, err)
	assert.Nil(t, sport)
}
