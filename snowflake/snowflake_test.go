package snowflake

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistributeId_NextId(t *testing.T) {
	id1, err := Default.NextId()
	assert.NoError(t, err, "snowflake next id error", err)
	assert.Greater(t, id1, int64(0), "snowflake next id less than 0")
}
