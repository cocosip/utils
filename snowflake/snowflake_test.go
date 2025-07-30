package snowflake

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mockable timeGen for testing clock rollback and sequence overflow
// We will modify the package-level timeGenFunc for testing purposes.

func TestNew(t *testing.T) {
	// Test valid creation
	id, err := New(DefaultEpoch, 5, 5, 12, 0, 0)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	// Test invalid workerId
	_, err = New(DefaultEpoch, 5, 5, 12, 32, 0) // max workerId is 31
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "worker ID can't be greater than")

	_, err = New(DefaultEpoch, 5, 5, 12, -1, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "worker ID can't be greater than")

	// Test invalid datacenterId
	_, err = New(DefaultEpoch, 5, 5, 12, 0, 32) // max datacenterId is 31
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "datacenter ID can't be greater than")

	_, err = New(DefaultEpoch, 5, 5, 12, 0, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "datacenter ID can't be greater than")
}

func TestDistributeId_NextId_Uniqueness(t *testing.T) {
	// Test uniqueness of generated IDs
	id, err := NewWithDefault(0, 0)
	assert.NoError(t, err)

	generatedIDs := make(map[int64]struct{})
	numIDs := 10000

	for i := 0; i < numIDs; i++ {
		newID, err := id.NextId()
		assert.NoError(t, err)
		_, exists := generatedIDs[newID]
		assert.False(t, exists, "Duplicate ID generated: %d", newID)
		generatedIDs[newID] = struct{}{}
	}
}

func TestDistributeId_NextId_Concurrency(t *testing.T) {
	// Test concurrency
	id, err := NewWithDefault(0, 0)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	numGoroutines := 10
	idsPerGoroutine := 1000
	generatedIDs := make(chan int64, numGoroutines*idsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				newID, err := id.NextId()
				assert.NoError(t, err)
				generatedIDs <- newID
			}
		}()
	}

	wg.Wait()
	close(generatedIDs)

	uniqueIDs := make(map[int64]struct{})
	for id := range generatedIDs {
		_, exists := uniqueIDs[id]
		assert.False(t, exists, "Duplicate ID generated in concurrent test: %d", id)
		uniqueIDs[id] = struct{}{}
	}
	assert.Equal(t, numGoroutines*idsPerGoroutine, len(uniqueIDs), "Not all IDs were unique")
}

func TestDistributeId_NextId_ClockRollback(t *testing.T) {
	// Test clock rollback
	originalTimeGenFunc := timeGenFunc
	defer func() {
		timeGenFunc = originalTimeGenFunc
	}()

	id, err := NewWithDefault(0, 0)
	assert.NoError(t, err)

	mockTime := int64(1000)
	timeGenFunc = func() int64 { return mockTime }

	_, err = id.NextId()
	assert.NoError(t, err)

	// Simulate clock moving backwards
	mockTime = 999
	_, err = id.NextId()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "clock moved backwards")
}

func TestDistributeId_NextId_SequenceOverflow(t *testing.T) {
	// Test sequence overflow
	originalTimeGenFunc := timeGenFunc
	defer func() {
		timeGenFunc = originalTimeGenFunc
	}()

	id, err := NewWithDefault(0, 0)
	assert.NoError(t, err)

	mockTime := int64(1000)
	timeGenFunc = func() int64 { return mockTime }

	id.lastTimestamp = 1000
	id.sequence = id.sequenceMask // Max sequence for current millisecond

	// Next ID should trigger sequence overflow and wait for next millisecond
	mockTime = 1000 // Keep time same to trigger overflow
	go func() {
		// Advance time slightly after a short delay to unblock tilNextMillis
		time.Sleep(10 * time.Millisecond)
		mockTime = 1001
	}()

	newID, err := id.NextId()
	assert.NoError(t, err)
	// Verify that the timestamp part of the new ID is from the next millisecond
	// (timestamp - epoch) << timestampLeftShift
	expectedTimestampPart := (1001 - id.epoch) << id.timestampLeftShift
	actualTimestampPart := newID & (^((1 << id.timestampLeftShift) - 1))
	assert.Equal(t, expectedTimestampPart, actualTimestampPart, "Timestamp part of ID mismatch after overflow")
}

func TestDistributeId_NextStringId(t *testing.T) {
	// Test NextStringId
	id, err := NewWithDefault(0, 0)
	assert.NoError(t, err)

	strID, err := id.NextStringId()
	assert.NoError(t, err)
	assert.NotEmpty(t, strID)

	// Try converting back to int64 to ensure it's a valid ID string
	intID, err := strconv.ParseInt(strID, 10, 64)
	assert.NoError(t, err)
	assert.Greater(t, intID, int64(0))
}

func TestNextId(t *testing.T) {
	// Test global NextId
	id, err := NextId()
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))
}

func TestNextStringId(t *testing.T) {
	// Test global NextStringId
	strID, err := NextStringId()
	assert.NoError(t, err)
	assert.NotEmpty(t, strID)
}
