package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_QueryFromMap(t *testing.T) {
	assert.Contains(t, getSportQueries()[sportsGet], "WHERE")
	assert.NotContains(t, getSportQueries()[sportsList], "WHERE")
}
