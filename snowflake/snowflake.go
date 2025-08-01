// Package snowflake implements the Snowflake ID generation algorithm.
// It generates unique, time-ordered, 64-bit IDs.
// The ID structure is: timestamp (41 bits) + datacenter ID (5 bits) + worker ID (5 bits) + sequence (12 bits).
package snowflake

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// DefaultEpoch is the default starting timestamp for Snowflake ID generation (2015-01-01 00:00:00 UTC).
	// All generated IDs will be relative to this epoch.
	DefaultEpoch int64 = 1420041600000 // 2015-01-01 00:00:00 UTC in milliseconds

	// DefaultWorkerIdBits defines the default number of bits allocated for the worker ID (5 bits).
	// This allows for 2^5 = 32 unique worker IDs (0-31).
	DefaultWorkerIdBits = 5

	// DefaultDatacenterIdBits defines the default number of bits allocated for the datacenter ID (5 bits).
	// This allows for 2^5 = 32 unique datacenter IDs (0-31).
	DefaultDatacenterIdBits = 5

	// DefaultSequenceBits defines the default number of bits allocated for the sequence number (12 bits).
	// This allows for 2^12 = 4096 IDs to be generated within a single millisecond.
	DefaultSequenceBits = 12

	// Default is a default singleton instance of DistributeId.
	// It is initialized with default epoch, worker ID 0, and datacenter ID 0.
	Default = newDefault()
)

// DistributeId represents a Snowflake ID generator instance.
// It holds the configuration and state required for ID generation.
type DistributeId struct {
	epoch              int64       // 开始时间截 (毫秒), ID 生成的基准时间
	workerIdBits       int         // 机器ID所占的位数 (通常为 5)
	datacenterIdBits   int         // 数据中心ID所占的位数 (通常为 5)
	sequenceBits       int         // 毫秒内序列号所占的位数 (通常为 12)
	workerId           int64       // 当前工作机器ID (0 ~ 2^workerIdBits - 1)
	datacenterId       int64       // 当前数据中心ID (0 ~ 2^datacenterIdBits - 1)
	workerIdShift      int         // 机器ID在 ID 中左移的位数 (sequenceBits)
	datacenterIdShift  int         // 数据中心ID在 ID 中左移的位数 (sequenceBits + workerIdBits)
	timestampLeftShift int         // 时间戳在 ID 中左移的位数 (datacenterIdShift + datacenterIdBits)
	sequenceMask       int64       // 毫秒内序列号的最大值掩码 (2^sequenceBits - 1)
	sequence           int64       // 毫秒内序列 (0 ~ sequenceMask)
	lastTimestamp      int64       // 上次生成ID的时间截 (毫秒)
	m                  *sync.Mutex // 互斥锁，用于保证并发安全
}

// newDefault initializes the default DistributeId instance.
// It panics if the default instance cannot be created, as it indicates a critical configuration error.
func newDefault() *DistributeId {
	id, err := NewWithDefault(0, 0)
	if err != nil {
		panic(fmt.Sprintf("Failed to create default snowflake instance: %v", err))
	}
	return id
}

// NewWithDefault creates a new DistributeId instance with default epoch and bit lengths.
// It takes workerId and datacenterId as parameters.
func NewWithDefault(workerId int64, datacenterId int64) (*DistributeId, error) {
	return New(DefaultEpoch, DefaultWorkerIdBits, DefaultDatacenterIdBits, DefaultSequenceBits, workerId, datacenterId)
}

// New creates a new DistributeId instance with custom configurations.
// It validates the workerId and datacenterId against their maximum allowed values.
func New(epoch int64, workerIdBits int, datacenterIdBits int, sequenceBits int, workerId int64, datacenterId int64) (*DistributeId, error) {
	// Calculate maximum allowed values for workerId and datacenterId based on their bit lengths.
	maxWorkerId := int64(-1 ^ (-1 << workerIdBits))
	maxDatacenterId := int64(-1 ^ (-1 << datacenterIdBits))

	// Validate workerId.
	if workerId > maxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker ID can't be greater than %d or less than 0", maxWorkerId)
	}
	// Validate datacenterId.
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenter ID can't be greater than %d or less than 0", maxDatacenterId)
	}

	// Initialize DistributeId instance.
	id := &DistributeId{
		epoch:            epoch,
		workerIdBits:     workerIdBits,
		datacenterIdBits: datacenterIdBits,
		sequenceBits:     sequenceBits,
		workerId:         workerId,
		datacenterId:     datacenterId,
		m:                &sync.Mutex{},
	}

	// Calculate bit shifts for combining components into a 64-bit ID.
	id.workerIdShift = sequenceBits
	id.datacenterIdShift = sequenceBits + workerIdBits
	id.timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits

	// Calculate the sequence mask.
	id.sequenceMask = -1 ^ (-1 << sequenceBits)

	return id, nil
}

// NextId generates a unique 64-bit Snowflake ID.
// It is thread-safe and handles clock rollback and sequence overflow.
func (d *DistributeId) NextId() (int64, error) {
	d.m.Lock()
	defer d.m.Unlock()

	timestamp := d.timeGen()

	// Handle clock rollback: if current timestamp is less than lastTimestamp, it means clock moved backwards.
	// This is a critical error as it can lead to duplicate IDs.
	if timestamp < d.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards. Refusing to generate ID for %d milliseconds", d.lastTimestamp-timestamp)
	}

	// If the current timestamp is the same as the last timestamp, increment the sequence.
	if d.lastTimestamp == timestamp {
		d.sequence = (d.sequence + 1) & d.sequenceMask
		// If sequence overflows (reaches sequenceMask + 1), wait for the next millisecond.
		if d.sequence == 0 {
			timestamp = d.tilNextMillis(d.lastTimestamp)
		}
	} else {
		// If the timestamp has changed, reset the sequence to 0.
		d.sequence = 0
	}

	// Update lastTimestamp for the next ID generation.
	d.lastTimestamp = timestamp

	// Combine all components into a 64-bit ID using bitwise operations:
	// (timestamp - epoch) << timestampLeftShift | (datacenterId << datacenterIdShift) | (workerId << workerIdShift) | sequence
	id := ((timestamp - d.epoch) << d.timestampLeftShift) |
		(d.datacenterId << d.datacenterIdShift) |
		(d.workerId << d.workerIdShift) |
		d.sequence

	return id, nil
}

// NextStringId generates a unique Snowflake ID and returns it as a string.
func (d *DistributeId) NextStringId() (string, error) {
	id, err := d.NextId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// timeGen returns the current timestamp in milliseconds since epoch.
// This function can be mocked for testing purposes.
var timeGenFunc = func() int64 {
	return time.Now().UnixMilli()
}

// timeGen returns the current timestamp in milliseconds since epoch.
// This function can be mocked for testing purposes.
func (d *DistributeId) timeGen() int64 {
	return timeGenFunc()
}

// tilNextMillis blocks until the next millisecond is reached.
// This is used when the sequence number for the current millisecond has been exhausted.
func (d *DistributeId) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGenFunc()
	for timestamp <= lastTimestamp {
		timestamp = timeGenFunc()
	}
	return timestamp
}

// NextId generates a unique 64-bit Snowflake ID using the default instance.
// This is a convenience function for quick ID generation without custom configuration.
func NextId() (int64, error) {
	return Default.NextId()
}

// NextStringId generates a unique Snowflake ID as a string using the default instance.
// This is a convenience function for quick ID generation without custom configuration.
func NextStringId() (string, error) {
	return Default.NextStringId()
}
