package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/colinSchofield/entain/racing/logging"
	"github.com/colinSchofield/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestOrderBy) ([]*racing.Race, error)

	// Get will return a single race or an sql.NotFound error.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, order *racing.ListRacesRequestOrderBy) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)
	query = r.applyOrderBy(query, order)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		wrappedError := fmt.Errorf("unexpected error in call to database Query: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// default or filter set to false, displays all races regardless of their visibility
	if filter.GetVisible() {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	logging.Logger().Debugf("Query: (%s) and args: (%v)", query, args)
	return query, args
}

// The order by must be from the 'approved' list of attribute names (default is 'advertisedStartTime').
func (r *racesRepo) applyOrderBy(query string, order *racing.ListRacesRequestOrderBy) string {

	if order == nil {
		return query
	}

	orderedQuery := fmt.Sprintf("%s ORDER BY %s %s", query, checkAttributeName(order.GetOrderBy()), order.GetDirection())
	return orderedQuery
}

// This is important as it protects against an sql injection attack. The default is to order by advertisedStartTime (i.e. 'advertised_start_time')
func checkAttributeName(attributeName string) string {

	const defaultOrderByColumn = "advertised_start_time"
	allowedColumnNames := map[string]string{
		"meetingId":           "meeting_id",
		"name":                "name",
		"visible":             "visible",
		"advertisedStartTime": "advertised_start_time",
	}

	if columnName, ok := allowedColumnNames[attributeName]; ok {
		return columnName
	} else {
		logging.Logger().Warnf("invalid attribute name of (%s). This may be a sql injection attack?!", attributeName)
		return defaultOrderByColumn
	}
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				logging.Logger().Debug("returning the empty element via error result of sql.ErrNoRows")
				return nil, nil
			}

			wrappedError := fmt.Errorf("unexpected error occurred in call to rows.Scan: %w", err)
			logging.Logger().Error(wrappedError)
			return nil, wrappedError
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			wrappedError := fmt.Errorf("unexpected error occurred in call to ptypes.TimestampProto: %w", err)
			logging.Logger().Error(wrappedError)
			return nil, wrappedError
		}

		race.AdvertisedStartTime = ts
		setRaceStatus(advertisedStart, &race)
		races = append(races, &race)
	}

	return races, nil
}

// Any race that is in the future is considered to be open
func setRaceStatus(advertisedStartTime time.Time, race *racing.Race) {

	if time.Now().Before(advertisedStartTime) {
		race.Status = racing.RaceStatus_OPEN
	} else {
		race.Status = racing.RaceStatus_CLOSED
	}
}

// Get will return a single race or an sql.NotFound error.
func (r *racesRepo) Get(id int64) (*racing.Race, error) {

	var (
		args []interface{}
	)
	query := getRaceQueries()[racesGet]
	args = append(args, id)

	if rows, err := r.db.Query(query, args...); err != nil {
		wrappedError := fmt.Errorf("unexpected error in call to database Query: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	} else if races, err := r.scanRaces(rows); err != nil {
		wrappedError := fmt.Errorf("unexpected error in call to database Query: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	} else if len(races) != 1 {
		return nil, sql.ErrNoRows
	} else {
		return races[0], nil
	}
}
