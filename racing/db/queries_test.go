package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_QueryFromMap(t *testing.T) {
	assert.Contains(t, getRaceQueries()[racesGet], "WHERE")
	assert.NotContains(t, getRaceQueries()[racesList], "WHERE")
}
