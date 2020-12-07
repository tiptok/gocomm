package sysx

import (
	"github.com/tiptok/gocomm/xstring"
	"os"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = xstring.RandId()
	}
}

func Hostname() string {
	return hostname
}
