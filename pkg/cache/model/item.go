package model

import (
	"encoding/json"
	"time"
)

type Item struct {
	Object     interface{} `json:"object"`     // object
	TTL        int         `json:"ttl"`        // key ttl, in second
	Outdate    int64       `json:"outdate"`    // outdated keys will be deleted from in-memory cache, but staty in redis.
	Expiration int64       `json:"expiration"` // expired keys will be deleted from redis.
	MarshData  []byte      `json:"-"`
}

// Returns true if data is outdated.
func (item Item) Expire() bool {
	if item.Outdate == 0 {
		return false
	}

	if item.Outdate < time.Now().UnixNano() {
		return true
	}
	return false
}

func (item Item) Data() []byte {
	return item.MarshData
}

func (item Item) String() string {
	d, _ := json.Marshal(item)
	return string(d)
}

func NewItem(v interface{}, d int) *Item {
	ttl := d
	var od, e int64
	if d > 0 {
		od = time.Now().Add(time.Duration(d) * time.Second).UnixNano()
		e = time.Now().Add(time.Duration(d*1) * time.Second).UnixNano() //lazyFactor
	}

	return &Item{
		Object:     v,
		TTL:        ttl,
		Outdate:    od,
		Expiration: e,
	}
}
