package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/colinSchofield/entain/sporting/logging"
	"github.com/colinSchofield/entain/sporting/proto/sporting"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of sports.
	List(filter *sporting.ListSportsRequestFilter, sort *sporting.ListSportsRequestOrderBy) ([]*sporting.Sport, error)

	// Get will return a single sport or an sql.NotFound error.
	Get(id int64) (*sporting.Sport, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport repository dummy data.
func (r *sportsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy sports.
		err = r.seed()
	})

	return err
}

func (r *sportsRepo) List(filter *sporting.ListSportsRequestFilter, order *sporting.ListSportsRequestOrderBy) ([]*sporting.Sport, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportQueries()[sportsList]

	query, args = r.applyFilter(query, filter)
	query = r.applyOrderBy(query, order)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		wrappedError := fmt.Errorf("unexpected error in call to database Query: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	}

	return r.scanSports(rows)
}

func (r *sportsRepo) applyFilter(query string, filter *sporting.ListSportsRequestFilter) (string, []interface{}) {
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

	// default or filter set to false, displays all sports regardless of their visibility
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
func (r *sportsRepo) applyOrderBy(query string, order *sporting.ListSportsRequestOrderBy) string {

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

func (m *sportsRepo) scanSports(
	rows *sql.Rows,
) ([]*sporting.Sport, error) {
	var sports []*sporting.Sport

	for rows.Next() {
		var sport sporting.Sport
		var advertisedStart time.Time

		if err := rows.Scan(&sport.Id, &sport.MeetingId, &sport.Name, &sport.Number, &sport.Visible, &advertisedStart); err != nil {
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

		sport.AdvertisedStartTime = ts
		setSportStatus(advertisedStart, &sport)
		sports = append(sports, &sport)
	}

	return sports, nil
}

// Any sport that is in the future is considered to be open
func setSportStatus(advertisedStartTime time.Time, sport *sporting.Sport) {

	if time.Now().Before(advertisedStartTime) {
		sport.Status = sporting.SportStatus_OPEN
	} else {
		sport.Status = sporting.SportStatus_CLOSED
	}
}

// Get will return a single sport or an sql.NotFound error.
func (r *sportsRepo) Get(id int64) (*sporting.Sport, error) {

	var (
		args []interface{}
	)
	query := getSportQueries()[sportsGet]
	args = append(args, id)

	if rows, err := r.db.Query(query, args...); err != nil {
		wrappedError := fmt.Errorf("unexpected error in call to database Query: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	} else if sports, err := r.scanSports(rows); err != nil {
		wrappedError := fmt.Errorf("unexpected error in call to database Query: %w", err)
		logging.Logger().Error(wrappedError)
		return nil, wrappedError
	} else if len(sports) != 1 {
		return nil, sql.ErrNoRows
	} else {
		return sports[0], nil
	}
}
