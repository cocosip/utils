package snowflake

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	DefaultEpoch            int64 = 1420041600000 //开始时间截(2015-01-01)
	DefaultWorkerIdBits           = 5             //机器id所占的位数(5)
	DefaultDatacenterIdBits       = 5             //数据标识id所占的位数(5)
	DefaultSequenceBits           = 12            //序列在id中占的位数(1ms内的并发数)(12)
	Default                       = newDefault()
)

type DistributeId struct {
	epoch              int64 //开始时间截(2015-01-01)
	workerIdBits       int   //机器id所占的位数
	datacenterIdBits   int   //数据标识id所占的位数
	sequenceBits       int   //序列在id中占的位数(1ms内的并发数)
	workerId           int64 //工作机器ID(0~31)
	datacenterId       int64 //数据中心ID(0~31)
	workerIdShift      int   //机器ID向左移12位
	datacenterIdShift  int   //数据标识id向左移17位(12+5)
	timestampLeftShift int   //时间截向左移22位(5+5+12)
	sequenceMask       int64 //生成序列的掩码，这里为4095 (0b111111111111=0xfff=4095)
	sequence           int64 //毫秒内序列(0~4095)
	lastTimestamp      int64 //上次生成ID的时间截
	m                  *sync.Mutex
}

func newDefault() *DistributeId {
	id, _ := NewWithDefault(0, 0)
	return id
}

func NewWithDefault(workerId int64, datacenterId int64) (*DistributeId, error) {
	return New(DefaultEpoch, DefaultWorkerIdBits, DefaultDatacenterIdBits, DefaultSequenceBits, workerId, datacenterId)
}

func New(epoch int64, workerIdBits int, datacenterIdBits int, sequenceBits int, workerId int64, datacenterId int64) (*DistributeId, error) {
	id := &DistributeId{
		epoch:            epoch,
		workerIdBits:     workerIdBits,
		datacenterIdBits: datacenterIdBits,
		sequenceBits:     sequenceBits,
		workerIdShift:    sequenceBits,
		workerId:         workerId,
		datacenterId:     datacenterId,
		m:                &sync.Mutex{},
	}

	id.datacenterIdShift = sequenceBits + workerIdBits
	id.timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	id.sequenceMask = -1 ^ (-1 << sequenceBits)

	maxWorkerId := int64(-1 ^ (-1 << workerIdBits))
	maxDatacenterId := int64(-1 ^ (-1 << datacenterIdBits))
	if workerId > maxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker Id can't be greater than %d or less than 0", maxWorkerId)
	}
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenter Id can't be greater than %d or less than 0", maxDatacenterId)
	}
	return id, nil
}

func (d *DistributeId) NextId() (int64, error) {
	d.m.Lock()
	defer d.m.Unlock()
	timestamp := d.timeGen()
	if timestamp < d.lastTimestamp {
		return 0, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %v milliseconds", d.lastTimestamp-timestamp)
	}

	if d.lastTimestamp == timestamp {
		d.sequence = (d.sequence + 1) & d.sequenceMask
		//毫秒内序列溢出
		if d.sequence == 0 {
			//阻塞到下一个毫秒,获得新的时间戳
			timestamp = d.tilNextMillis(d.lastTimestamp)
		}
	} else {
		//时间戳改变，毫秒内序列重置
		d.sequence = 0
	}

	//上次生成ID的时间截
	d.lastTimestamp = timestamp

	//移位并通过或运算拼到一起组成64位的ID
	id := ((timestamp - d.epoch) << d.timestampLeftShift) | (d.datacenterId << d.datacenterIdShift) | (d.workerId << d.workerIdShift) | d.sequence
	return id, nil
}

func (d *DistributeId) NextStringId() (string, error) {
	id, err := d.NextId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (d *DistributeId) timeGen() int64 {
	return time.Now().UnixMilli()
}

// tilNextMillis 阻塞到下一个毫秒，直到获得新的时间戳
func (d *DistributeId) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := d.timeGen()
	for {
		if timestamp <= lastTimestamp {
			timestamp = d.timeGen()
		} else {
			break
		}
	}
	return timestamp

}
