package idgen

import (
	"fmt"
)

var gsf *Sonyflake // github.com/sony/sonyflake v1.0.0

func init() {
	st := Settings{
		MachineID: getMachineId,
	}
	gsf = NewSonyflake(st)
}

func getMachineId() (uint16, error) {
	// TODO
	return uint16(1), nil
}

func Init(midGenFunc func() (uint16, error)) {
	st := Settings{
		MachineID: midGenFunc,
	}
	gsf = NewSonyflake(st)
}

// Next generates next id as an uint64
func Next() (id int64) {
	var i uint64
	if gsf != nil {
		i, _ = gsf.NextID()
		id = int64(i)
	}
	return
}

// NextString generates next id as a string
func NextString() (id string) {
	id = fmt.Sprintf("%d", Next())

	return
}

func GetOne() int64 {
	return Next()
}

func GetMulti(n int) (ids []int64) {
	for i := 0; i < n; i++ {
		ids = append(ids, Next())
	}
	return
}
